package logic

import (
	"time"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/api/model"
	"golang.org/x/net/context"
)

func init() {
	RegisterTasker(DiffTasker, diffTasker)
}

func diffTasker(ctx context.Context, result model.KeyedResult, task model.KeyedTask, user model.User, startTime time.Time) (skills model.Skills, err error) {
	var submissions model.Submissions
	var submissionKeys []*datastore.Key
	submissionKeys, err = model.NewQueryForSubmission().
		Ancestor(result.Key).
		Filter("Task =", task.Key).
		Order("Time").
		GetAll(ctx, &submissions)
	if err != nil {
		return
	}

	var cs float64
	if len(submissions) > 0 {
		cs, err = diffCodingSpeed(ctx, submissions.Key(submissionKeys), task, startTime)
		if err != nil {
			return
		}
	}

	skills.CodingSpeed = cs
	return
}

func diffCodingSpeed(ctx context.Context, submissions []model.KeyedSubmission, task model.KeyedTask, startTime time.Time) (cs float64, err error) {
	userCodingTime := submissions[len(submissions)-1].Time.Sub(startTime)
	insDel, err := getInsertedDeleted(submissions)
	if err != nil {
		return 0, err
	}

	// TODO(victorbalan, flowlo): Load the testResults for the last built submission.
	var testResults []model.DiffTestResult
	_, err = model.NewQueryForDiffTestResult().
		Ancestor(submissions[len(submissions)-1].Key).
		Order("Start").
		GetAll(ctx, &testResults)
	if err != nil {
		return
	}

	_, red := diffRedGreenTests(testResults)
	if red != 0 {
		return 0, nil
	}
	return codingSpeedValue(len(submissions), len(testResults),
		userCodingTime, task.Assignment.Duration,
		insDel.Inserted, insDel.Deleted,
		0.4, 0.3, 0.3)
}

func diffRedGreenTests(testResults []model.DiffTestResult) (green, red int) {
	for _, tr := range testResults {
		if tr.DiffLines == nil {
			green++
		} else {
			red++
		}
	}
	return
}
