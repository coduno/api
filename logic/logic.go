package logic

import (
	"strconv"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

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
type TaskerFunc func(ctx context.Context, task, user *datastore.Key) (model.Skills, error)

const (
	// Average computes the weighted average over all task results. It is included
	// in this package.
	Average Resulter = 1 + iota
	maxResulter
)

const (
	maxTasker Tasker = iota
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
func (t Tasker) Call(ctx context.Context, task, user *datastore.Key) (model.Skills, error) {
	if t > 0 && t < maxTasker {
		f := taskers[t]
		if f != nil {
			return f(ctx, task, user)
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
			taskResult, err := Tasker(task.Tasker).Call(ctx, challenge.Tasks[i], resultKey.Parent().Parent())
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
}
