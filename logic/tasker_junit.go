package logic

import (
	"github.com/coduno/api/model"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func init() {
	RegisterTasker(JunitTasker, junitTasker)
}

func junitTasker(ctx context.Context, taskKey, resultKey, userKey *datastore.Key) (skills model.Skills, err error) {
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

	var cs float64
	cs, err = codingSpeed(submissions, task, result.StartTimes[getTaskIndex(challenge, taskKey)])
	if err != nil {
		return
	}

	skills.CodingSpeed = cs
	return
}
