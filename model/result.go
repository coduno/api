package model

import (
	"time"

	"google.golang.org/appengine/datastore"
)

//go:generate generator -c "Submission"

// Result holds the performance of an User for some
// Challenge. It is fixed once the user has completed
// the Challenge or timed out. After that, only
// Skills are to be modified in case the
// logic in Challenge changes.
//
// Saved in Datastore, Result will be a child
// entity to Profile, so keys pointing to a Result
// can be used to obtain the Profile they influence.
type Result struct {
	// Calculated by logic from the Challenge. If
	// it is missing an average over all FinalSubmissions
	// will be computed at best effort.
	Skills `datastore:",index",json:",omitempty"`

	// Challenge refers to the challenge that this
	// result provides data for.
	Challenge *datastore.Key `datastore:",index",json:",omitempty"`

	// Indicates when the user has started to work on
	// a Task (meaning as soon as the Task
	// is served to the user).
	//
	// In case all Tasks are available to the
	// user in parallel, it is possible that every
	// element of this slice holds the same value.
	// Anyway, the Challenge logic has to make sense
	// of this property and how to interpret it.
	//
	// Indexed the same as Challenge.Tasks.
	StartTimes []time.Time `datastore:",noindex",json:",omitempty"`

	// Points to the last submission to the
	// corresponding Task.
	//
	// Indexed the same as Challenge.Tasks.
	FinalSubmissions []*datastore.Key `datastore:",noindex",json:",omitempty"`

	// The time when the coder started the Challenge.
	Started time.Time `datastore:",index",json:",omitempty"`

	// The time when the coder finished the Challenge.
	Finished time.Time `datastore:",index",json:",omitempty"`

	// When this result was last (re)computed by the
	// Resulter in Challenge.
	Computed time.Time `datastore:",index",json:",omitempty"`
}
