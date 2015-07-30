package model

import "time"

// Assignment conveys the specification of what a
// User must do in order to fulfill a Task or
// Challenge.
type Assignment struct {
	// Name should give a very brief and memorable
	// description and classification of the assignment.
	Name string

	// Description should detail what this assignment is
	// about and why it makes sense.
	//
	// TODO(victorbalan): Improve documentation.
	Description string

	// Instructions should make clear how the assignment
	// is to be carried out (e.g. a step-by-step guide).
	//
	// NOTE(flowlo, victorbalan): Instructions is not guaranteed
	// to be backwards-compatible. In the future it may
	// be an URL pointing at the source of the instructions.
	Instructions string

	// Expected time to complete the assignment. Can serve
	// as deadline.
	Duration time.Duration

	// Where to deliver results as part of carrying out the
	// assignment.
	Endpoints Endpoints
}
