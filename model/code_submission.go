package model

import "time"

//go:generate generator

type CodeSubmission struct {
	Submission

	Code,
	Language,
	Stdout,
	Stderr,
	Prepare string

	Rusage     Rusage
	Exit       error
	Start, End time.Time
}
