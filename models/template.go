package models

import "google.golang.org/appengine/datastore"

// Template contains data about a code tamplate assigned to a challenge
type Template struct {
	EntityID  string         `json:"id"`
	Language  string         `json:"language"`
	Path      string         `json:"path"`
	Challenge *datastore.Key `json:"challenge"`
}
