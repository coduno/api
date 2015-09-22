package logic

import (
	"errors"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/coduno/api/model"
)

// Resulter is a reference to a unique implementation of a ResulterFunc.
type Resulter int

// ResulterFunc is a function that will compute skills for the referenced
// Result. It will call Taskers of encapsulated tasks on demand.
type ResulterFunc func(ctx context.Context, result model.KeyedResult, challenge model.Challenge) error

// Tasker is a reference to a unique implementation of a TaskerFunc.
type Tasker int

// TaskerFunc is a function that will compute task results for the given Task
// and User.
type TaskerFunc func(ctx context.Context, result model.KeyedResult, task model.KeyedTask, user model.User, startTime time.Time) (model.Skills, error)

const (
	// Average computes the weighted average over all task results. It is included
	// in this package.
	Average Resulter = 1 + iota
	maxResulter
)

const (
	// JunitTasker computes the skills for a specific task. It iterates
	// over all the submissions.
	JunitTasker Tasker = 1 + iota
	DiffTasker
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
func (r Resulter) Call(ctx context.Context, result model.KeyedResult, challenge model.Challenge) error {
	if r > 0 && r < maxResulter {
		f := resulters[r]
		if f != nil {
			return f(ctx, result, challenge)
		}
	}
	panic("logic: requested resulter function #" + strconv.Itoa(int(r)) + " is unavailable")
}

// Call looks up a registered Tasker and calls it.
func (t Tasker) Call(ctx context.Context, result model.KeyedResult, task model.KeyedTask, user model.User, startTime time.Time) (model.Skills, error) {
	if t > 0 && t < maxTasker {
		f := taskers[t]
		if f != nil {
			return f(ctx, result, task, user, startTime)
		}
	}
	return model.Skills{}, errors.New("logic: requested tasker function #" + strconv.Itoa(int(t)) + " is unavailable")
}
