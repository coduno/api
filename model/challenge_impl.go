// This file was automatically generated from
//
//	challenge.go
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

// ChallengeKind is the kind used in Datastore to store entities Challenge entities.
const ChallengeKind = "Challenge"

// Challenges is just a slice of Challenge.
type Challenges []Challenge

// KeyedChallenge is a struct that embeds Challenge and also contains a Key, mainly used for encoding to JSON.
type KeyedChallenge struct {
	*Challenge
	Key *datastore.Key
}

// Key is a shorthand to fill a KeyedChallenge with an entity and it's key.
func (ƨ *Challenge) Key(key *datastore.Key) *KeyedChallenge {
	return &KeyedChallenge{
		Challenge: ƨ,
		Key:       key,
	}
}

// Key is a shorthand to fill a slice of KeyedChallenge with some entities alongside their keys.
func (ƨ Challenges) Key(keys []*datastore.Key) (keyed []KeyedChallenge) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedChallenge, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedChallenge{
			Challenge: &ƨ[i],
			Key:       keys[i],
		}
	}
	return
}

// Put will put this Challenge into Datastore using the given key.
func (ƨ Challenge) Put(ctx context.Context, key *datastore.Key) (*datastore.Key, error) {
	if key != nil {
		return datastore.Put(ctx, key, &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Challenge", nil), &ƨ)
}

// PutWithParent can be used to save this Challenge as child of another
// entity.
// This will error if parent == nil.
func (ƨ Challenge) PutWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Challenge", parent), &ƨ)
}

// NewQueryForChallenge prepares a datastore.Query that can be
// used to query entities of type Challenge.
func NewQueryForChallenge() *datastore.Query {
	return datastore.NewQuery("Challenge")
}
