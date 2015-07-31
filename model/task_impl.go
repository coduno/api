// This file was automatically generated from
//
//	task.go
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

type Tasks []Task

type KeyedTask struct {
	Task *Task
	Key  *datastore.Key
}

func (ƨ *Task) Key(key *datastore.Key) *KeyedTask {
	return &KeyedTask{
		Task: ƨ,
		Key:  key,
	}
}

func (ƨ Tasks) Key(keys []*datastore.Key) (keyed []KeyedTask) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedTask, 0, len(ƨ))
	for i := range keyed {
		keyed = append(keyed, KeyedTask{
			Task: &ƨ[i],
			Key:  keys[i],
		})
	}
	return
}
