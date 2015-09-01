// This file was automatically generated from
//
//	submission.go
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

// SubmissionKind is the kind used in Datastore to store entities Submission entities.
const SubmissionKind = "Submission"

// Submissions is just a slice of Submission.
type Submissions []Submission

// KeyedSubmission is a struct that embeds Submission and also contains a Key, mainly used for encoding to JSON.
type KeyedSubmission struct {
	*Submission
	Key *datastore.Key
}

// Key is a shorthand to fill a KeyedSubmission with an entity and it's key.
func (ƨ *Submission) Key(key *datastore.Key) *KeyedSubmission {
	return &KeyedSubmission{
		Submission: ƨ,
		Key:        key,
	}
}

// Key is a shorthand to fill a slice of KeyedSubmission with some entities alongside their keys.
func (ƨ Submissions) Key(keys []*datastore.Key) (keyed []KeyedSubmission) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedSubmission, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedSubmission{
			Submission: &ƨ[i],
			Key:        keys[i],
		}
	}
	return
}

// Put will put this Submission into Datastore using the given key.
func (ƨ Submission) Put(ctx context.Context, key *datastore.Key) (*datastore.Key, error) {
	if key != nil {
		return datastore.Put(ctx, key, &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Submission", nil), &ƨ)
}

// PutWithParent can be used to save this Submission as child of another
// entity.
// This will error if parent == nil.
func (ƨ Submission) PutWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Submission", parent), &ƨ)
}

// NewQueryForSubmission prepares a datastore.Query that can be
// used to query entities of type Submission.
func NewQueryForSubmission() *datastore.Query {
	return datastore.NewQuery("Submission")
}
