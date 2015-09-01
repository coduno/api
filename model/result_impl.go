// This file was automatically generated from
//
//	result.go
//
// by
//
//	generator -c Submission
//
// DO NOT EDIT

package model

import (
	"errors"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// ResultKind is the kind used in Datastore to store entities Result entities.
const ResultKind = "Result"

// Results is just a slice of Result.
type Results []Result

// KeyedResult is a struct that embeds Result and also contains a Key, mainly used for encoding to JSON.
type KeyedResult struct {
	*Result
	Key *datastore.Key
}

// Key is a shorthand to fill a KeyedResult with an entity and it's key.
func (ƨ *Result) Key(key *datastore.Key) *KeyedResult {
	return &KeyedResult{
		Result: ƨ,
		Key:    key,
	}
}

// Key is a shorthand to fill a slice of KeyedResult with some entities alongside their keys.
func (ƨ Results) Key(keys []*datastore.Key) (keyed []KeyedResult) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedResult, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedResult{
			Result: &ƨ[i],
			Key:    keys[i],
		}
	}
	return
}

// Put will put this Result into Datastore using the given key.
func (ƨ Result) Put(ctx context.Context, key *datastore.Key) (*datastore.Key, error) {
	if key != nil {
		return datastore.Put(ctx, key, &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Result", nil), &ƨ)
}

// PutWithParent can be used to save this Result as child of another
// entity.
// This will error if parent == nil.
func (ƨ Result) PutWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Result", parent), &ƨ)
}

// NewQueryForResult prepares a datastore.Query that can be
// used to query entities of type Result.
func NewQueryForResult() *datastore.Query {
	return datastore.NewQuery("Result")
}
