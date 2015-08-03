// This file was automatically generated from
//
//	invitation.go
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

type Invitations []Invitation

type KeyedInvitation struct {
	*Invitation
	Key *datastore.Key
}

func (ƨ *Invitation) Key(key *datastore.Key) *KeyedInvitation {
	return &KeyedInvitation{
		Invitation: ƨ,
		Key:        key,
	}
}

func (ƨ Invitations) Key(keys []*datastore.Key) (keyed []KeyedInvitation) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedInvitation, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedInvitation{
			Invitation: &ƨ[i],
			Key:        keys[i],
		}
	}
	return
}
