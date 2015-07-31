// This file was automatically generated from
//
//	result.go
//
// by
//
//	generator -c Submission
//
// DO NOT EDIT

package model

import (
	"google.golang.org/appengine/datastore"
)

type Results []Result

type KeyedResult struct {
	Result *Result
	Key    *datastore.Key
}

func (ƨ *Result) Key(key *datastore.Key) *KeyedResult {
	return &KeyedResult{
		Result: ƨ,
		Key:    key,
	}
}

func (ƨ Results) Key(keys []*datastore.Key) (keyed []KeyedResult) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedResult, 0, len(ƨ))
	for i := range keyed {
		keyed = append(keyed, KeyedResult{
			Result: &ƨ[i],
			Key:    keys[i],
		})
	}
	return
}
