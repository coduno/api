// This file was automatically generated from
//
//	code_task.go
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

type CodeTasks []CodeTask

type KeyedCodeTask struct{
	CodeTask *CodeTask
	Key      *datastore.Key
}

func (ƨ *CodeTask) Key(key *datastore.Key) (*KeyedCodeTask) {
	return &KeyedCodeTask{
		CodeTask: ƨ,
		Key:      key,
	}
}

func (ƨ CodeTasks) Key(keys []*datastore.Key) (keyed []KeyedCodeTask) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedCodeTask, 0, len(ƨ))
	for i := range keyed {
		keyed = append(keyed, KeyedCodeTask{
			CodeTask: &ƨ[i],
			Key:      keys[i],
		})
	}
	return
}
