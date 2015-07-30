// This file was automatically generated from
//
//	task.go
//
// by
//
//	generator
//
// at
//
//	2015-07-30T17:21:42+03:00
//
// Do not edit it!

package model

import (
	"errors"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// Save will put this Task into Datastore using the given key.
func (ƨ Task) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "tasks", nil), &ƨ)
}

// SaveWithParent can be used to save this Task as child of another
// entity.
// This will error if parent == nil.
func (ƨ Task) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "tasks", parent), &ƨ)
}

// NewQueryForTask prepares a datastore.Query that can be
// used to query entities of type Task.
func NewQueryForTask() *datastore.Query {
	return datastore.NewQuery("tasks")
}

type TaskHandler struct{}

func (ƨ TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.URL.Path == "" {
		var results Tasks
		keys, _ := NewQueryForTask().GetAll(ctx, &results)
		results.Write(w, keys)
		return
	}

	k, _ := datastore.DecodeKey(r.URL.Path)
	var ƨ_ Task
	datastore.Get(ctx, k, &ƨ_)
	ƨ_.Write(w, k)
}

func ServeTask(prefix string, muxes ...*http.ServeMux) {
	path := prefix + "tasks" + "/"

	if len(muxes) == 0 {
		http.Handle(path, http.StripPrefix(path, TaskHandler{}))
	}

	for _, mux := range muxes {
		mux.Handle(path, http.StripPrefix(path, TaskHandler{}))
	}
}
