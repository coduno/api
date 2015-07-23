package models

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

//CoderKind the name of the collection in datastore
const CoderKind = "coders"

// Coder contains the data related to a coder
type Coder struct {
	EntityID  int64  `datastore:"-" json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Save a new coder
func (coder Coder) Save(ctx context.Context) (*datastore.Key, error) {
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, CoderKind, nil), &coder)
}
