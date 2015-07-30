package model

import (
	"time"

	"google.golang.org/appengine/datastore"
)

//go:generate generator

type Invitation struct {
	User *datastore.Key
	Sent time.Time
}
