// This file was automatically generated from
//
//	diff_test_result.go
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

// DiffTestResultKind is the kind used in Datastore to store entities DiffTestResult entities.
const DiffTestResultKind = "DiffTestResult"

// DiffTestResults is just a slice of DiffTestResult.
type DiffTestResults []DiffTestResult

// KeyedDiffTestResult is a struct that embeds DiffTestResult and also contains a Key, mainly used for encoding to JSON.
type KeyedDiffTestResult struct {
	*DiffTestResult
	Key *datastore.Key
}

// Key is a shorthand to fill a KeyedDiffTestResult with an entity and it's key.
func (ƨ *DiffTestResult) Key(key *datastore.Key) *KeyedDiffTestResult {
	return &KeyedDiffTestResult{
		DiffTestResult: ƨ,
		Key:            key,
	}
}

// Key is a shorthand to fill a slice of KeyedDiffTestResult with some entities alongside their keys.
func (ƨ DiffTestResults) Key(keys []*datastore.Key) (keyed []KeyedDiffTestResult) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedDiffTestResult, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedDiffTestResult{
			DiffTestResult: &ƨ[i],
			Key:            keys[i],
		}
	}
	return
}

// Put will put this DiffTestResult into Datastore using the given key.
func (ƨ DiffTestResult) Put(ctx context.Context, key *datastore.Key) (*datastore.Key, error) {
	if key != nil {
		return datastore.Put(ctx, key, &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "DiffTestResult", nil), &ƨ)
}

// PutWithParent can be used to save this DiffTestResult as child of another
// entity.
// This will error if parent == nil.
func (ƨ DiffTestResult) PutWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "DiffTestResult", parent), &ƨ)
}

// NewQueryForDiffTestResult prepares a datastore.Query that can be
// used to query entities of type DiffTestResult.
func NewQueryForDiffTestResult() *datastore.Query {
	return datastore.NewQuery("DiffTestResult")
}
