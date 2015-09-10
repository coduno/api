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

func averageResulter(ctx context.Context, resultKey *datastore.Key) error {
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

	_, err := result.Put(ctx, resultKey)
	return err
}
