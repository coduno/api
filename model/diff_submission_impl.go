// This file was automatically generated from
//
//	diff_submission.go
//
// by
//
//	generator
//
// DO NOT EDIT

package model

import (
	"errors"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// DiffSubmissionKind is the kind used in Datastore to store entities of type DiffSubmission.
const DiffSubmissionKind = "DiffSubmission"

// DiffSubmissions is just a slice of DiffSubmission.
type DiffSubmissions []DiffSubmission

// KeyedDiffSubmission is a struct that embeds DiffSubmission and also contains a Key, mainly used for encoding to JSON.
type KeyedDiffSubmission struct {
	*DiffSubmission
	Key *datastore.Key
}

// Key is a shorthand to fill a KeyedDiffSubmission with an entity and it's key.
func (ƨ *DiffSubmission) Key(key *datastore.Key) *KeyedDiffSubmission {
	return &KeyedDiffSubmission{
		DiffSubmission: ƨ,
		Key:            key,
	}
}

// Key is a shorthand to fill a slice of KeyedDiffSubmission with some entities alongside their keys.
func (ƨ DiffSubmissions) Key(keys []*datastore.Key) (keyed []KeyedDiffSubmission) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedDiffSubmission, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedDiffSubmission{
			DiffSubmission: &ƨ[i],
			Key:            keys[i],
		}
	}
	return
}

// Save will put this DiffSubmission into Datastore using the given key.
func (ƨ DiffSubmission) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "DiffSubmission", nil), &ƨ)
}

// SaveWithParent can be used to save this DiffSubmission as child of another
// entity.
// This will error if parent == nil.
func (ƨ DiffSubmission) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "DiffSubmission", parent), &ƨ)
}

// NewQueryForDiffSubmission prepares a datastore.Query that can be
// used to query entities of type DiffSubmission.
func NewQueryForDiffSubmission() *datastore.Query {
	return datastore.NewQuery("DiffSubmission")
}
