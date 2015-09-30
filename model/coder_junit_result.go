package model

//go:generate generator

// CoderJunitTestResult is the result of a set of JUnit tests written by the
// coder for some specs ran against code written by the coduno team for said
// specs. Besides the results of the test, it also encapsulates wether it was
// supposed to fail in order to assess the correctness of the unit tests.
type CoderJunitTestResult struct {
	JunitTestResult
	ShouldFail bool
}
