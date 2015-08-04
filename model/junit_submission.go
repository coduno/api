package model

import "time"

//go:generate generator

type JunitSubmission struct {
	Submission

	Code,
	Stdout,
	Stderr,
	Exit string

	Start, End time.Time

	Result UnitTestResult
}

// UnitTestResult holds the unit test result created by JUnit.
type UnitTestResult struct {
	Tests    int        `xml:"tests,attr"`
	Failures int        `xml:"failures,attr"`
	Errors   int        `xml:"errors,attr"`
	TestCase []TestCase `xml:"testcase"`
}

// TestCase holds a test case created by JUnit
type TestCase struct {
	Name     string        `xml:"name,attr"`
	Duration time.Duration `xml:"time,attr"`
	Failure  Failure       `xml:"failure"`
}

// Failure holds a failure created by JUnit.
type Failure struct {
	Message string `xml:"message,attr"`
}