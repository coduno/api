package models

import "google.golang.org/appengine/datastore"

// Challenge contains the data of a challenge with the company that created it
type Challenge struct {
	EntityID     string         `json:"id"`
	Name         string         `json:"name"`
	Instructions string         `json:"instructions"`
	Company      *datastore.Key `json:"company"`
}
