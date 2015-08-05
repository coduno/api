// This file was automatically generated from
//
//	profile.go
//
// by
//
//	generator -c Result
//
// DO NOT EDIT

package model

import (
	"errors"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// ProfileKind is the kind used in Datastore to store entities of type Profile.
const ProfileKind = "Profile"

// Profiles is just a slice of Profile.
type Profiles []Profile

// KeyedProfile is a struct that embeds Profile and also contains a Key, mainly used for encoding to JSON.
type KeyedProfile struct {
	*Profile
	Key *datastore.Key
}

// Key is a shorthand to fill a KeyedProfile with an entity and it's key.
func (ƨ *Profile) Key(key *datastore.Key) *KeyedProfile {
	return &KeyedProfile{
		Profile: ƨ,
		Key:     key,
	}
}

// Key is a shorthand to fill a slice of KeyedProfile with some entities alongside their keys.
func (ƨ Profiles) Key(keys []*datastore.Key) (keyed []KeyedProfile) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedProfile, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedProfile{
			Profile: &ƨ[i],
			Key:     keys[i],
		}
	}
	return
}

// Save will put this Profile into Datastore using the given key.
func (ƨ Profile) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Profile", nil), &ƨ)
}

// SaveWithParent can be used to save this Profile as child of another
// entity.
// This will error if parent == nil.
func (ƨ Profile) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Profile", parent), &ƨ)
}

// NewQueryForProfile prepares a datastore.Query that can be
// used to query entities of type Profile.
func NewQueryForProfile() *datastore.Query {
	return datastore.NewQuery("Profile")
}
