package logic

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"time"

	"github.com/coduno/api/model"
)

func codingSpeed(submissions []model.Submission, task model.Task, startTime time.Time) (cs float64, err error) {
	// TODO(victorbalan): Load it from the params map
	nrOfTests := 5
	userCodingTime := submissions[len(submissions)-1].Time.Sub(startTime)

	var oldCode io.Reader
	oldCode, err = readFromGCS(submissions[0].Code)
	if err != nil {
		return
	}
	var initialCode []byte
	if initialCode, err = ioutil.ReadAll(oldCode); err != nil {
		return
	}
	insDel := &InsDel{len(bytes.Split(initialCode, []byte("\n"))), 0}
	// Iterate all submissions
	for i := 1; i < len(submissions); i++ {
		var newCode io.Reader
		newCode, err = readFromGCS(submissions[i].Code)
		if err != nil {
			return
		}
		var id InsDel
		if id, err = computeInsertedDeletedLines(newCode, oldCode); err != nil {
			return
		}
		insDel.Add(id)

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

	return codingSpeedValue(len(submissions), nrOfTests,
		userCodingTime, task.Assignment.Duration,
		insDel.Inserted, insDel.Deleted,
		0.4, 0.3, 0.3)
}

func codingSpeedValue(userSubmissions, nrOfTests int,
	userCodingTime, maxCodingTime time.Duration,
	insertedLines, deletedLines int,
	nosWeight, timetWeight, idlWeight float64) (float64, error) {

	if nosWeight+timetWeight+idlWeight != 1 {
		return 0, errors.New("weights do not sum up to 1")
	}
	return nOfSubmissionsTest(userSubmissions, nrOfTests)*nosWeight +
		timeTest(userCodingTime, maxCodingTime)*timetWeight +
		insertedDeletedLinesTest(insertedLines, deletedLines)*idlWeight, nil
}

// TODO(victorbalan): This test is deprecated with the new structure. Needs
// refactoring.
//
// if x<0 return insanity
// x is between [0..1]
// s is the number of submissions the user made
// t is the total number of tests for the task(eg. 5 JUnit tests, 3 diff tests)
func nOfSubmissionsTest(s, t int) (x float64) {
	if s == 1 {
		return 1
	}
	x = 1 - float64(s)/(float64(t)*10)
	if x < 0 {
		x = 0
	}
	return
}

// x is between [0..1]
// d is the time it took the coder to complete the task
// tmax is the maximum allowed time for the task
func timeTest(d, tmax time.Duration) float64 {
	return 1 - d.Seconds()/(tmax.Seconds()/2)
}

// this is the initial version of this test
// the future versions will take in the differences from task to task
// of the nr of inserted/deleted lines together with the correct answers diff
//
// x is between [0..1]
// i is the number of inserted lines
// d is tthe number of deleted lines
func insertedDeletedLinesTest(i, d int) float64 {
	return 1 - float64(d)/float64(i)
}
