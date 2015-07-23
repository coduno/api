package models

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

//ChallangeKind name of the collection in datastore
const ChallangeKind = "challange"

// Challenge contains the data of a challenge with the company that created it
type Challenge struct {
	EntityID     string         `json:"id"`
	Name         string         `json:"name"`
	Instructions string         `json:"instructions"`
	Company      *datastore.Key `json:"company"`
}

//Save a new challagne to the datastore
func (challange Challenge) Save(ctx context.Context) (*datastore.Key, error) {
	key, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, ChallangeKind, nil), &challange)
	if err != nil {
		return nil, err
	}
	return key, nil
}
