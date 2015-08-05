package model

import "time"

//go:generate generator

// CodeSubmission represents a submission of some piece of code.
// Besides the submitted code it also contains output generate by
// compilation and/or execution, resource usage etc.
type CodeSubmission struct {
	Submission

	Code,
	Language,
	Stdout,
	Stderr,
	Exit,
	Prepare string

	Rusage     Rusage
	Start, End time.Time
}
