package logic

import (
	"errors"
	"time"
)

func codingSpeed(userSubmissions, nrOfTests int,
	userCodingTime, maxCodingTime time.Duration,
	insertedLines, deletedLines int,
	nosWeight, timetWeight, idlWeight float64) (float64, error) {

	if nosWeight+timetWeight+idlWeight != 1 {
		return 0, errors.New("Weights do not sum up to 0.")
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
