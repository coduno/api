// This file was automatically generated from
//
//	coder_junit_result.go
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

// CoderJunitTestResultKind is the kind used in Datastore to store entities CoderJunitTestResult entities.
const CoderJunitTestResultKind = "CoderJunitTestResult"

// CoderJunitTestResults is just a slice of CoderJunitTestResult.
type CoderJunitTestResults []CoderJunitTestResult

// KeyedCoderJunitTestResult is a struct that embeds CoderJunitTestResult and also contains a Key, mainly used for encoding to JSON.
type KeyedCoderJunitTestResult struct {
	*CoderJunitTestResult
	Key *datastore.Key
}

// Key is a shorthand to fill a KeyedCoderJunitTestResult with an entity and it's key.
func (ƨ *CoderJunitTestResult) Key(key *datastore.Key) *KeyedCoderJunitTestResult {
	return &KeyedCoderJunitTestResult{
		CoderJunitTestResult: ƨ,
		Key:                  key,
	}
}

// Key is a shorthand to fill a slice of KeyedCoderJunitTestResult with some entities alongside their keys.
func (ƨ CoderJunitTestResults) Key(keys []*datastore.Key) (keyed []KeyedCoderJunitTestResult) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedCoderJunitTestResult, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedCoderJunitTestResult{
			CoderJunitTestResult: &ƨ[i],
			Key:                  keys[i],
		}
	}
	return
}

// Put will put this CoderJunitTestResult into Datastore using the given key.
func (ƨ CoderJunitTestResult) Put(ctx context.Context, key *datastore.Key) (*datastore.Key, error) {
	if key != nil {
		return datastore.Put(ctx, key, &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "CoderJunitTestResult", nil), &ƨ)
}

// PutWithParent can be used to save this CoderJunitTestResult as child of another
// entity.
// This will error if parent == nil.
func (ƨ CoderJunitTestResult) PutWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "CoderJunitTestResult", parent), &ƨ)
}

// NewQueryForCoderJunitTestResult prepares a datastore.Query that can be
// used to query entities of type CoderJunitTestResult.
func NewQueryForCoderJunitTestResult() *datastore.Query {
	return datastore.NewQuery("CoderJunitTestResult")
}
