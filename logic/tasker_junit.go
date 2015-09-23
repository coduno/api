package logic

import (
	"time"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/api/model"
	"golang.org/x/net/context"
)

func init() {
	RegisterTasker(JunitTasker, junitTasker)
}

func junitTasker(ctx context.Context, result model.KeyedResult, task model.KeyedTask, user model.User, startTime time.Time) (skills model.Skills, err error) {
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
		cs, err = junitCodingSpeed(ctx, submissions.Key(submissionKeys), task, startTime)
		if err != nil {
			return
		}
	}

	skills.CodingSpeed = cs
	return
}

func junitCodingSpeed(ctx context.Context, submissions []model.KeyedSubmission, task model.KeyedTask, startTime time.Time) (cs float64, err error) {
	userCodingTime := submissions[len(submissions)-1].Time.Sub(startTime)
	insDel, err := getInsertedDeleted(submissions)
	if err != nil {
		return 0, err
	}

	// TODO(victorbalan, flowlo): Take in account the nr of green/red tests
	// We load only the test results for the last submission now.
	// Load the testResults for the last built submission.
	var testResults []model.JunitTestResult
	_, err = model.NewQueryForJunitTestResult().
		Ancestor(submissions[len(submissions)-1].Key).
		Order("Start").
		GetAll(ctx, &testResults)
	if err != nil {
		return
	}

	nrOfTests, _, red := junitRedGreenTests(testResults)
	if red != 0 {
		return 0, nil
	}
	return codingSpeedValue(len(submissions), nrOfTests,
		userCodingTime, task.Assignment.Duration,
		insDel.Inserted, insDel.Deleted,
		0.4, 0.3, 0.3)
}

func junitRedGreenTests(testResults []model.JunitTestResult) (nrOfTests, green, red int) {
	for _, tr := range testResults {
		nrOfTests += len(tr.Results.TestCase)
		for _, r := range tr.Results.TestCase {
			if r.Failure.Message == "" {
				green++
			} else {
				red++
			}
		}
	}
	return
}
