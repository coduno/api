package model

import (
	"google.golang.org/appengine/datastore"
)

//go:generate generator

// Challenge is an abstract piece of work that can consist
// of many different Tasks.
//
// Saved in Datastore, Challenge will be a child
// entity to Company, so keys pointing to a Challenge
// can be used to obtain the Company that owns it.
type Challenge struct {
	Assignment

	// The tasks that have to be fulfilled in order
	// to successfully complete the Challenge.
	//
	// Result.StartTimes and Result.FinalSubmissions
	// depend on the ordering of this slice. Also it
	// affects the rendering of this Challenge with
	// respect to the user. Therefore it must be
	// guaranteed to be stable.
	Tasks []*datastore.Key

	// To normalize the results of all Tasks
	//
	// TODO(flowlo): Clear specification.
	Logic logic
}
