// This file was automatically generated from
//
//	code_submission.go
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

type CodeSubmissions []CodeSubmission

type KeyedCodeSubmission struct {
	*CodeSubmission
	Key *datastore.Key
}

func (ƨ *CodeSubmission) Key(key *datastore.Key) *KeyedCodeSubmission {
	return &KeyedCodeSubmission{
		CodeSubmission: ƨ,
		Key:            key,
	}
}

func (ƨ CodeSubmissions) Key(keys []*datastore.Key) (keyed []KeyedCodeSubmission) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedCodeSubmission, 0, len(ƨ))
	for i := range keyed {
		keyed = append(keyed, KeyedCodeSubmission{
			CodeSubmission: &ƨ[i],
			Key:            keys[i],
		})
	}
	return
}
