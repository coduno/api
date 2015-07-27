package models

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

//FingerprintKind name of the collection in datastore
const FingerprintKind = "fingerprints"

// Fingerprint contains data that links a coder to a challenge
type Fingerprint struct {
	Coder     *datastore.Key `json:"coder"`
	Challenge *datastore.Key `json:"challenge"`
	Token     string         `json:"token"`
}

// Save a new fingerprint
func (fingerprint Fingerprint) Save(ctx context.Context) (*datastore.Key, error) {
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, FingerprintKind, nil), &fingerprint)
}
