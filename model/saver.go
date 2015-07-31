package model

import (
	"encoding/json"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Saver means every model that can be saved to Datastore
// by calling the uniform Save function.
type Saver interface {
	Save(ctx context.Context) (key *datastore.Key, err error)
}

type Keyer interface {
	Key(key *datastore.Key) json.Marshaler
}
