package model

import "time"

//go:generate generator -c "Result"

// Profile is the current rating of a User. It
// can vary over times as new Results come in
// and will be recomputed as needed.
// It is there to give a quick overview over the
// skills/performance of an individual user.
//
// For short-timed challenges, the Profile should
// not be updated, but only when a final result
// was produced.
//
// Long-term challenges may decide to refresh the
// competing user's Profile as pleased, but should
// not do so more than once a day.
//
// Saved in Datastore, Profile will be a child
// entity to User, so keys pointing to a Profile
// can be used to obtain the user they represent.
type Profile struct {
	Skills
	LastUpdate time.Time
}
