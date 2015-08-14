package model

import "time"

//go:generate generator

// JunitSubmission is a submission to a set of JUnit test cases.
// Besides the uploaded code, it also encapsulates more detailed
// results generated by JUnit.
type JunitSubmission struct {
	Submission

	Code   string `datastore:",noindex"`
	Stdout string `datastore:",noindex"`
	Stderr string `datastore:",noindex"`
	Exit   string `datastore:",noindex"`

	Start time.Time `datastore:",index"`
	End   time.Time `datastore:",index"`

	Results UnitTestResults `datastore:",noindex"`
}

// UnitTestResults holds the unit test result created by JUnit.
type UnitTestResults struct {
	Tests    int        `xml:"tests,attr"`
	Failures int        `xml:"failures,attr"`
	Errors   int        `xml:"errors,attr"`
	TestCase []TestCase `xml:"testcase"`
}

// TestCase holds a test case created by JUnit.
type TestCase struct {
	Name     string        `xml:"name,attr"`
	Duration time.Duration `xml:"time,attr"`
	Failure  Failure       `xml:"failure"`
}

// Failure holds a failure created by JUnit.
type Failure struct {
	Message string `xml:"message,attr"`
}
