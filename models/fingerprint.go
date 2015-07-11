package models

import "google.golang.org/appengine/datastore"

// Fingerprint -
type Fingerprint struct {
	EntityID  string         `json:"id"`
	Coder     *datastore.Key `json:"coder"`
	Challenge *datastore.Key `json:"challenge"`
	Token     string         `json:"token"`
}
