// This file was automatically generated from
//
//	robot_test_results.go
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

// RobotTestResultsKind is the kind used in Datastore to store entities RobotTestResults entities.
const RobotTestResultsKind = "RobotTestResults"

// RobotTestResultss is just a slice of RobotTestResults.
type RobotTestResultss []RobotTestResults

// KeyedRobotTestResults is a struct that embeds RobotTestResults and also contains a Key, mainly used for encoding to JSON.
type KeyedRobotTestResults struct {
	*RobotTestResults
	Key *datastore.Key
}

// Key is a shorthand to fill a KeyedRobotTestResults with an entity and it's key.
func (ƨ *RobotTestResults) Key(key *datastore.Key) *KeyedRobotTestResults {
	return &KeyedRobotTestResults{
		RobotTestResults: ƨ,
		Key:              key,
	}
}

// Key is a shorthand to fill a slice of KeyedRobotTestResults with some entities alongside their keys.
func (ƨ RobotTestResultss) Key(keys []*datastore.Key) (keyed []KeyedRobotTestResults) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedRobotTestResults, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedRobotTestResults{
			RobotTestResults: &ƨ[i],
			Key:              keys[i],
		}
	}
	return
}

// Put will put this RobotTestResults into Datastore using the given key.
func (ƨ RobotTestResults) Put(ctx context.Context, key *datastore.Key) (*datastore.Key, error) {
	if key != nil {
		return datastore.Put(ctx, key, &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "RobotTestResults", nil), &ƨ)
}

// PutWithParent can be used to save this RobotTestResults as child of another
// entity.
// This will error if parent == nil.
func (ƨ RobotTestResults) PutWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "RobotTestResults", parent), &ƨ)
}

// NewQueryForRobotTestResults prepares a datastore.Query that can be
// used to query entities of type RobotTestResults.
func NewQueryForRobotTestResults() *datastore.Query {
	return datastore.NewQuery("RobotTestResults")
}
