// This file was automatically generated from
//
//	submission.go
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

type Submissions []Submission

type KeyedSubmission struct{
	Submission *Submission
	Key      *datastore.Key
}

func (ƨ *Submission) Key(key *datastore.Key) (*KeyedSubmission) {
	return &KeyedSubmission{
		Submission: ƨ,
		Key:      key,
	}
}

func (ƨ Submissions) Key(keys []*datastore.Key) (keyed []KeyedSubmission) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedSubmission, 0, len(ƨ))
	for i := range keyed {
		keyed = append(keyed, KeyedSubmission{
			Submission: &ƨ[i],
			Key:      keys[i],
		})
	}
	return
}
