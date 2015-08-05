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
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// TaskKind is the kind used in Datastore to store entities of type Task.
const TaskKind = "Task"

// Tasks is just a slice of Task.
type Tasks []Task

// KeyedTask is a struct that embeds Task and also contains a Key, mainly used for encoding to JSON.
type KeyedTask struct {
	*Task
	Key *datastore.Key
}

// Key is a shorthand to fill a KeyedTask with an entity and it's key.
func (ƨ *Task) Key(key *datastore.Key) *KeyedTask {
	return &KeyedTask{
		Task: ƨ,
		Key:  key,
	}
}

// Key is a shorthand to fill a slice of KeyedTask with some entities alongside their keys.
func (ƨ Tasks) Key(keys []*datastore.Key) (keyed []KeyedTask) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedTask, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedTask{
			Task: &ƨ[i],
			Key:  keys[i],
		}
	}
	return
}

// Save will put this Task into Datastore using the given key.
func (ƨ Task) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Task", nil), &ƨ)
}

// SaveWithParent can be used to save this Task as child of another
// entity.
// This will error if parent == nil.
func (ƨ Task) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Task", parent), &ƨ)
}

// NewQueryForTask prepares a datastore.Query that can be
// used to query entities of type Task.
func NewQueryForTask() *datastore.Query {
	return datastore.NewQuery("Task")
}
