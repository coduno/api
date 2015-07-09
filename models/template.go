package models

import "google.golang.org/appengine/datastore"

// Template -
type Template struct {
	EntityID  string         `json:"id"`
	Language  string         `json:"language"`
	Path      string         `json:"path"`
	Challenge *datastore.Key `json:"challenge"`
}
