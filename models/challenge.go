package models

import "google.golang.org/appengine/datastore"

// Challenge -
type Challenge struct {
	EntityID     string         `json:"id"`
	Name         string         `json:"name"`
	Instructions string         `json:"instructions"`
	Company      *datastore.Key `json:"company"`
}
