package logic

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/cloud/storage"

	"github.com/coduno/api/model"
)

// Resulter is a reference to a unique implementation of a ResulterFunc.
type Resulter int

// ResulterFunc is a function that will compute skills for the referenced
// Result. It will call Taskers of encapsulated tasks on demand.
type ResulterFunc func(ctx context.Context, resultKey *datastore.Key) error

// Tasker is a reference to a unique implementation of a TaskerFunc.
type Tasker int

// TaskerFunc is a function that will compute task results for the given Task
// and User.
type TaskerFunc func(ctx context.Context, task, resultKey, user *datastore.Key) (model.Skills, error)

const (
	// Average computes the weighted average over all task results. It is included
	// in this package.
	Average Resulter = 1 + iota
	maxResulter
)

const (
	JunitTasker Tasker = 1 + iota
	maxTasker
)

var resulters = make([]ResulterFunc, maxResulter)
var taskers = make([]TaskerFunc, maxTasker)

// RegisterResulter registers a function to be called lated via Call. It is
// usually called from the init function of a package that contains a Resulter.
func RegisterResulter(r Resulter, f ResulterFunc) {
	if r >= maxResulter {
		panic("logic: RegisterResulter of unknown resulter function")
	}
	resulters[r] = f
}

// RegisterTasker registers a function to be called lated via Call. It is
// usually called from the init function of a package that contains a Tasker.
func RegisterTasker(t Tasker, f TaskerFunc) {
	if t >= maxTasker {
		panic("logic: RegisterTasker of unknown tasker function")
	}
	taskers[t] = f
}

// Call looks up a registered Resulter and calls it.
func (r Resulter) Call(ctx context.Context, resultKey *datastore.Key) error {
	if r > 0 && r < maxResulter {
		f := resulters[r]
		if f != nil {
			return f(ctx, resultKey)
		}
	}
	panic("logic: requested resulter function #" + strconv.Itoa(int(r)) + " is unavailable")
}

// Call looks up a registered Tasker and calls it.
func (t Tasker) Call(ctx context.Context, task, result, user *datastore.Key) (model.Skills, error) {
	if t > 0 && t < maxTasker {
		f := taskers[t]
		if f != nil {
			return f(ctx, task, result, user)
		}
	}
	panic("logic: requested tasker function #" + strconv.Itoa(int(t)) + " is unavailable")
}

func init() {
	RegisterResulter(Average, func(ctx context.Context, resultKey *datastore.Key) error {
		var result model.Result
		if err := datastore.Get(ctx, resultKey, &result); err != nil {
			return err
		}

		var challenge model.Challenge
		if err := datastore.Get(ctx, result.Challenge, &challenge); err != nil {
			return err
		}

		var tasks model.Tasks
		if err := datastore.GetMulti(ctx, challenge.Tasks, &tasks); err != nil {
			return err
		}

		weightSum := model.Skills{} // this could be SkillWeights, but would need more conversions
		average := model.Skills{}

		for i, task := range tasks {
			taskResult, err := Tasker(task.Tasker).Call(ctx, challenge.Tasks[i], resultKey, resultKey.Parent().Parent())
			if err != nil {
				return err
			}
			average = average.Add(taskResult.Mul(model.Skills(task.SkillWeights)))
			weightSum = weightSum.Add(model.Skills(task.SkillWeights))
		}

		result.Skills = average.Div(weightSum)
		result.Computed = time.Now()

		_, err := result.Save(ctx, resultKey)
		return err
	})

	RegisterTasker(JunitTasker, func(ctx context.Context, taskKey, resultKey, userKey *datastore.Key) (skills model.Skills, err error) {
		var submissions []model.Submission
		_, err = model.NewQueryForSubmission().
			Ancestor(resultKey).
			Filter("Task =", taskKey).
			Order("Start").
			GetAll(ctx, submissions)
		if err != nil {
			return
		}
		var task model.Task
		if err = datastore.Get(ctx, taskKey, &task); err != nil {
			return
		}

		var result model.Result
		if err = datastore.Get(ctx, resultKey, &result); err != nil {
			return
		}

		var challenge model.Challenge
		if err = datastore.Get(ctx, result.Challenge, &challenge); err != nil {
			return
		}

		// TODO(victorbalan): Load it from the params map
		nrOfTests := 5
		userCodingTime := submissions[len(submissions)].Time.Sub(result.StartTimes[getTaskIndex(challenge, taskKey)])

		splitByNewLine := func(c rune) bool {
			return c == '\n'
		}

		var oldCode string
		oldCode, err = loadFromGCS(ctx, submissions[0].Code)
		if err != nil {
			return
		}

		insDel := &InsDel{len(strings.FieldsFunc(oldCode, splitByNewLine)), 0}
		// Iterate all submissions
		for i := 1; i < len(submissions); i++ {
			var newCode string
			newCode, err = loadFromGCS(ctx, submissions[i].Code)
			if err != nil {
				return
			}
			insDel.Add(computeInsertedDeletedLines(newCode, oldCode, splitByNewLine))

			// TODO(victorbalan, flowlo): Take in account the nr of green/red tests
			// var testResultKeys []*datastore.Key
			// testResultKeys, err = model.NewQueryForJunitTestResult().
			// 	Ancestor(submissionKeys[i]).
			// 	Order("Start").
			// 	GetAll(ctx, nil)
			// if err != nil {
			// 	return
			// }
			oldCode = newCode
		}

		var cs float64
		cs, err = codingSpeed(len(submissions), nrOfTests,
			userCodingTime, task.Assignment.Duration,
			insDel.Inserted, insDel.Deleted,
			0.4, 0.3, 0.3)

		skills.CodingSpeed = cs
		return
	})
}

// InsDel holds the number of inserted and deleted lines
type InsDel struct {
	Inserted,
	Deleted int
}

func (id *InsDel) Add(insDel InsDel) {
	id.Inserted += insDel.Inserted
	id.Deleted += insDel.Deleted
}

func computeInsertedDeletedLines(oldCode, newCode string, splitFunc func(c rune) bool) InsDel {
	var i, d int
	currentFields := strings.FieldsFunc(newCode, splitFunc)
	oldFields := strings.FieldsFunc(oldCode, splitFunc)
	for _, val := range currentFields {
		if !strings.Contains(oldCode, val) {
			i++
		}
	}
	for _, val := range oldFields {
		if !strings.Contains(oldCode, val) {
			d++
		}
	}
	return InsDel{i, d}
}

func getTaskIndex(c model.Challenge, task *datastore.Key) int {
	for i, val := range c.Tasks {
		if val.Equal(task) {
			return i
		}
	}
	return -1
}

func loadFromGCS(ctx context.Context, so model.StoredObject) (string, error) {
	rc, err := storage.NewReader(ctx, so.Bucket, so.Name)
	if err != nil {
		return "", err
	}
	defer rc.Close()
	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
