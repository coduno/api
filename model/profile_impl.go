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
	"google.golang.org/appengine/datastore"
)

type Profiles []Profile

type KeyedProfile struct {
	*Profile
	Key *datastore.Key
}

func (ƨ *Profile) Key(key *datastore.Key) *KeyedProfile {
	return &KeyedProfile{
		Profile: ƨ,
		Key:     key,
	}
}

func (ƨ Profiles) Key(keys []*datastore.Key) (keyed []KeyedProfile) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedProfile, 0, len(ƨ))
	for i := range keyed {
		keyed = append(keyed, KeyedProfile{
			Profile: &ƨ[i],
			Key:     keys[i],
		})
	}
	return
}
