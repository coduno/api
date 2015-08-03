package model

import "time"

//go:generate generator

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
