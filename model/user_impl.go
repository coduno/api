// This file was automatically generated from
//
//	user.go
//
// by
//
//	generator -c Profile
//
// DO NOT EDIT

package model

import (
	"google.golang.org/appengine/datastore"
)

type Users []User

type KeyedUser struct {
	User *User
	Key  *datastore.Key
}

func (ƨ *User) Key(key *datastore.Key) *KeyedUser {
	return &KeyedUser{
		User: ƨ,
		Key:  key,
	}
}

func (ƨ Users) Key(keys []*datastore.Key) (keyed []KeyedUser) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedUser, 0, len(ƨ))
	for i := range keyed {
		keyed = append(keyed, KeyedUser{
			User: &ƨ[i],
			Key:  keys[i],
		})
	}
	return
}
