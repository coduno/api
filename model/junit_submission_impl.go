// This file was automatically generated from
//
//	junit_submission.go
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

// JunitSubmissionKind is the kind used in Datastore to store entities of type JunitSubmission.
const JunitSubmissionKind = "JunitSubmission"

// JunitSubmissions is just a slice of JunitSubmission.
type JunitSubmissions []JunitSubmission

// KeyedJunitSubmission is a struct that embeds JunitSubmission and also contains a Key, mainly used for encoding to JSON.
type KeyedJunitSubmission struct {
	*JunitSubmission
	Key *datastore.Key
}

// Key is a shorthand to fill a KeyedJunitSubmission with an entity and it's key.
func (ƨ *JunitSubmission) Key(key *datastore.Key) *KeyedJunitSubmission {
	return &KeyedJunitSubmission{
		JunitSubmission: ƨ,
		Key:             key,
	}
}

// Key is a shorthand to fill a slice of KeyedJunitSubmission with some entities alongside their keys.
func (ƨ JunitSubmissions) Key(keys []*datastore.Key) (keyed []KeyedJunitSubmission) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedJunitSubmission, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedJunitSubmission{
			JunitSubmission: &ƨ[i],
			Key:             keys[i],
		}
	}
	return
}

// Save will put this JunitSubmission into Datastore using the given key.
func (ƨ JunitSubmission) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "JunitSubmission", nil), &ƨ)
}

// SaveWithParent can be used to save this JunitSubmission as child of another
// entity.
// This will error if parent == nil.
func (ƨ JunitSubmission) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "JunitSubmission", parent), &ƨ)
}

// NewQueryForJunitSubmission prepares a datastore.Query that can be
// used to query entities of type JunitSubmission.
func NewQueryForJunitSubmission() *datastore.Query {
	return datastore.NewQuery("JunitSubmission")
}
