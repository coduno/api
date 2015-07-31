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
	"google.golang.org/appengine/datastore"
)

type Challenges []Challenge

type KeyedChallenge struct{
	Challenge *Challenge
	Key      *datastore.Key
}

func (ƨ *Challenge) Key(key *datastore.Key) (*KeyedChallenge) {
	return &KeyedChallenge{
		Challenge: ƨ,
		Key:      key,
	}
}

func (ƨ Challenges) Key(keys []*datastore.Key) (keyed []KeyedChallenge) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedChallenge, 0, len(ƨ))
	for i := range keyed {
		keyed = append(keyed, KeyedChallenge{
			Challenge: &ƨ[i],
			Key:      keys[i],
		})
	}
	return
}
