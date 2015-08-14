// This file was automatically generated from
//
//	code_submission.go
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

// CodeSubmissionKind is the kind used in Datastore to store entities CodeSubmission entities.
const CodeSubmissionKind = "CodeSubmission"

// CodeSubmissions is just a slice of CodeSubmission.
type CodeSubmissions []CodeSubmission

// KeyedCodeSubmission is a struct that embeds CodeSubmission and also contains a Key, mainly used for encoding to JSON.
type KeyedCodeSubmission struct {
	*CodeSubmission
	Key *datastore.Key
}

// Key is a shorthand to fill a KeyedCodeSubmission with an entity and it's key.
func (ƨ *CodeSubmission) Key(key *datastore.Key) *KeyedCodeSubmission {
	return &KeyedCodeSubmission{
		CodeSubmission: ƨ,
		Key:            key,
	}
}

// Key is a shorthand to fill a slice of KeyedCodeSubmission with some entities alongside their keys.
func (ƨ CodeSubmissions) Key(keys []*datastore.Key) (keyed []KeyedCodeSubmission) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedCodeSubmission, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedCodeSubmission{
			CodeSubmission: &ƨ[i],
			Key:            keys[i],
		}
	}
	return
}

// Save will put this CodeSubmission into Datastore using the given key.
func (ƨ CodeSubmission) Save(ctx context.Context, key *datastore.Key) (*datastore.Key, error) {
	if key != nil {
		return datastore.Put(ctx, key, &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "CodeSubmission", nil), &ƨ)
}

// SaveWithParent can be used to save this CodeSubmission as child of another
// entity.
// This will error if parent == nil.
func (ƨ CodeSubmission) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "CodeSubmission", parent), &ƨ)
}

// NewQueryForCodeSubmission prepares a datastore.Query that can be
// used to query entities of type CodeSubmission.
func NewQueryForCodeSubmission() *datastore.Query {
	return datastore.NewQuery("CodeSubmission")
}
