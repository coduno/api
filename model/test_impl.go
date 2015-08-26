// This file was automatically generated from
//
//	test.go
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

// TestKind is the kind used in Datastore to store entities Test entities.
const TestKind = "Test"

// Tests is just a slice of Test.
type Tests []Test

// KeyedTest is a struct that embeds Test and also contains a Key, mainly used for encoding to JSON.
type KeyedTest struct {
	*Test
	Key *datastore.Key
}

// Key is a shorthand to fill a KeyedTest with an entity and it's key.
func (ƨ *Test) Key(key *datastore.Key) *KeyedTest {
	return &KeyedTest{
		Test: ƨ,
		Key:  key,
	}
}

// Key is a shorthand to fill a slice of KeyedTest with some entities alongside their keys.
func (ƨ Tests) Key(keys []*datastore.Key) (keyed []KeyedTest) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedTest, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedTest{
			Test: &ƨ[i],
			Key:  keys[i],
		}
	}
	return
}

// Save will put this Test into Datastore using the given key.
func (ƨ Test) Save(ctx context.Context, key *datastore.Key) (*datastore.Key, error) {
	if key != nil {
		return datastore.Put(ctx, key, &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Test", nil), &ƨ)
}

// SaveWithParent can be used to save this Test as child of another
// entity.
// This will error if parent == nil.
func (ƨ Test) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Test", parent), &ƨ)
}

// NewQueryForTest prepares a datastore.Query that can be
// used to query entities of type Test.
func NewQueryForTest() *datastore.Query {
	return datastore.NewQuery("Test")
}
