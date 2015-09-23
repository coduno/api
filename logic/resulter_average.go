package logic

import (
	"time"

	"github.com/coduno/api/model"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func init() {
	RegisterResulter(Average, averageResulter)
}

func averageResulter(ctx context.Context, result model.KeyedResult, challenge model.Challenge) error {
	tasks := make([]model.Task, len(challenge.Tasks))
	if err := datastore.GetMulti(ctx, challenge.Tasks, tasks); err != nil {
		return err
	}

	var user model.User
	if err := datastore.Get(ctx, result.Key.Parent().Parent(), &user); err != nil {
		return err
	}

	var nrOfComputations float64
	average := model.Skills{}

	for i, task := range tasks {
		taskResult, err := Tasker(task.Tasker).Call(ctx, result, *task.Key(challenge.Tasks[i]), user, result.StartTimes[getTaskIndex(challenge, challenge.Tasks[i])])
		if err != nil {
			// TODO: ignore error for now. We`ll treat it after we have all the taskers available
			//return err
		} else {
			average = average.Add(taskResult)
			nrOfComputations++
		}
	}

	result.Skills = average.DivBy(nrOfComputations)
	result.Computed = time.Now()

	_, err := result.Put(ctx, result.Key)
	return err
}
