package model

import "google.golang.org/appengine/datastore"

//go:generate generator

// Template contains data about a code template assigned to a Task
type Template struct {
	Language  string         `datastore:",index"`
	Path      string         `datastore:",noindex"`
	Challenge *datastore.Key `datastore:",index"`
}
