package model

import (
	"time"

	"google.golang.org/appengine/datastore"
)

//go:generate generator

// Invitation represents the message sent by a company to a User
// in order to ask them to do a Challenge.
type Invitation struct {
	User *datastore.Key `datastore:",index"`
	Sent time.Time      `datastore:",index"`
}
