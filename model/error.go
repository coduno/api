package model

import (
	"time"

	"google.golang.org/appengine/datastore"
)

//go:generate generator

// Error contains the data of an
// error encountered by the user.
type Error struct {
	time.Time
	Status int
	User   *datastore.Key
	// The description of the events that
	// led to the error as provided by the user.
	Description string
	// The route to the page where
	// the error was encountered.
	Route string
}
