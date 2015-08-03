// This file was automatically generated from
//
//	junit_submission.go
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

type JunitSubmissions []JunitSubmission

type KeyedJunitSubmission struct {
	*JunitSubmission
	Key *datastore.Key
}

func (ƨ *JunitSubmission) Key(key *datastore.Key) *KeyedJunitSubmission {
	return &KeyedJunitSubmission{
		JunitSubmission: ƨ,
		Key:             key,
	}
}

func (ƨ JunitSubmissions) Key(keys []*datastore.Key) (keyed []KeyedJunitSubmission) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedJunitSubmission, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedJunitSubmission{
			JunitSubmission: &ƨ[i],
			Key:             keys[i],
		}
	}
	return
}
