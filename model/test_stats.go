package model

import "google.golang.org/appengine/datastore"

type TestStats struct {
	Stdout,
	Stderr string
	Test   *datastore.Key
	Failed bool
}
