package model

import "time"

//go:generate generator

// CodeSubmission represents a submission of some piece of code.
// Besides the submitted code it also contains output generate by
// compilation and/or execution, resource usage etc.
type CodeSubmission struct {
	Submission

	Code     string `datastore:",noindex"`
	Language string `datastore:",index"`
	Stdout   string `datastore:",noindex"`
	Stderr   string `datastore:",noindex"`
	Exit     string `datastore:",noindex"`
	Prepare  string `datastore:",noindex"`

	Rusage Rusage    `datastore:",noindex"`
	Start  time.Time `datastore:",index"`
	End    time.Time `datastore:",index"`
}
