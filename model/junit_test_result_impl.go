// This file was automatically generated from
//
//	junit_test_result.go
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

// JunitTestResultKind is the kind used in Datastore to store entities JunitTestResult entities.
const JunitTestResultKind = "JunitTestResult"

// JunitTestResults is just a slice of JunitTestResult.
type JunitTestResults []JunitTestResult

// KeyedJunitTestResult is a struct that embeds JunitTestResult and also contains a Key, mainly used for encoding to JSON.
type KeyedJunitTestResult struct {
	*JunitTestResult
	Key *datastore.Key
}

// Key is a shorthand to fill a KeyedJunitTestResult with an entity and it's key.
func (ƨ *JunitTestResult) Key(key *datastore.Key) *KeyedJunitTestResult {
	return &KeyedJunitTestResult{
		JunitTestResult: ƨ,
		Key:             key,
	}
}

// Key is a shorthand to fill a slice of KeyedJunitTestResult with some entities alongside their keys.
func (ƨ JunitTestResults) Key(keys []*datastore.Key) (keyed []KeyedJunitTestResult) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedJunitTestResult, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedJunitTestResult{
			JunitTestResult: &ƨ[i],
			Key:             keys[i],
		}
	}
	return
}

// Save will put this JunitTestResult into Datastore using the given key.
func (ƨ JunitTestResult) Save(ctx context.Context, key *datastore.Key) (*datastore.Key, error) {
	if key != nil {
		return datastore.Put(ctx, key, &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "JunitTestResult", nil), &ƨ)
}

// SaveWithParent can be used to save this JunitTestResult as child of another
// entity.
// This will error if parent == nil.
func (ƨ JunitTestResult) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "JunitTestResult", parent), &ƨ)
}

// NewQueryForJunitTestResult prepares a datastore.Query that can be
// used to query entities of type JunitTestResult.
func NewQueryForJunitTestResult() *datastore.Query {
	return datastore.NewQuery("JunitTestResult")
}
