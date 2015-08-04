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
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type CodeTasks []CodeTask

type KeyedCodeTask struct {
	*CodeTask
	Key *datastore.Key
}

func (ƨ *CodeTask) Key(key *datastore.Key) *KeyedCodeTask {
	return &KeyedCodeTask{
		CodeTask: ƨ,
		Key:      key,
	}
}

func (ƨ CodeTasks) Key(keys []*datastore.Key) (keyed []KeyedCodeTask) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedCodeTask, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedCodeTask{
			CodeTask: &ƨ[i],
			Key:      keys[i],
		}
	}
	return
}

// Save will put this CodeTask into Datastore using the given key.
func (ƨ CodeTask) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "CodeTask", nil), &ƨ)
}

// SaveWithParent can be used to save this CodeTask as child of another
// entity.
// This will error if parent == nil.
func (ƨ CodeTask) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "CodeTask", parent), &ƨ)
}

// NewQueryForCodeTask prepares a datastore.Query that can be
// used to query entities of type CodeTask.
func NewQueryForCodeTask() *datastore.Query {
	return datastore.NewQuery("CodeTask")
}

type CodeTaskHandler struct{}

func (ƨ CodeTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.URL.Path == "" {
		var results CodeTasks
		keys, _ := NewQueryForCodeTask().GetAll(ctx, &results)
		json.NewEncoder(w).Encode(results.Key(keys))
		return
	}

	k, _ := datastore.DecodeKey(r.URL.Path)
	var entity CodeTask
	datastore.Get(ctx, k, &entity)
	json.NewEncoder(w).Encode(entity)
}

func ServeCodeTask(prefix string, muxes ...*http.ServeMux) {
	path := prefix + "CodeTask" + "/"

	if len(muxes) == 0 {
		http.Handle(path, http.StripPrefix(path, CodeTaskHandler{}))
	}

	for _, mux := range muxes {
		mux.Handle(path, http.StripPrefix(path, CodeTaskHandler{}))
	}
}
