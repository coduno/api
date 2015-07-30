package model

import (
	"time"

	"google.golang.org/appengine/datastore"
)

//go:generate generator

// Submission is a form of result for some
// Task.
//
// This type is very general and should be embedded in more
// concrete types, so that it matches the expectations of the
// Task it corresponds to. For example:
//
//	type QuizSubmission struct {
//		Submission
//
//		CorrectAnswers,
//		TotalAnswers int
//	}
//
//	type CodeSubmission struct {
//		Submission
//
//		Code,
//		Language string
//	}
//
// TODO(flowlo): Switching parent key from Task to Result
// after completion of Task? Does that make sense?
//
//	if key.Parent().Kind() == "task" {
//		// yay, task is still running
//	} else if key.Parent().Kind() == "result" {
//		// we can do fast lookup of the submission now, no problem
//	} else {
//		// error
//	}
type Submission struct {
	Time time.Time
	Task *datastore.Key
}

// /api/users/bfsdlf/profiles/jkadfgsdfg/results/sldfng/codeSubmissions/lasdbfd/
