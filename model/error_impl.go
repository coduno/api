// This file was automatically generated from
//
//	error.go
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

type Errors []Error

type KeyedError struct {
	*Error
	Key *datastore.Key
}

func (ƨ *Error) Key(key *datastore.Key) *KeyedError {
	return &KeyedError{
		Error: ƨ,
		Key:   key,
	}
}

func (ƨ Errors) Key(keys []*datastore.Key) (keyed []KeyedError) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedError, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedError{
			Error: &ƨ[i],
			Key:   keys[i],
		}
	}
	return
}

// Save will put this Error into Datastore using the given key.
func (ƨ Error) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Error", nil), &ƨ)
}

// SaveWithParent can be used to save this Error as child of another
// entity.
// This will error if parent == nil.
func (ƨ Error) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Error", parent), &ƨ)
}

// NewQueryForError prepares a datastore.Query that can be
// used to query entities of type Error.
func NewQueryForError() *datastore.Query {
	return datastore.NewQuery("Error")
}

type ErrorHandler struct{}

func (ƨ ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.URL.Path == "" {
		var results Errors
		keys, _ := NewQueryForError().GetAll(ctx, &results)
		json.NewEncoder(w).Encode(results.Key(keys))
		return
	}

	k, _ := datastore.DecodeKey(r.URL.Path)
	var entity Error
	datastore.Get(ctx, k, &entity)
	json.NewEncoder(w).Encode(entity)
}

func ServeError(prefix string, muxes ...*http.ServeMux) {
	path := prefix + "Error" + "/"

	if len(muxes) == 0 {
		http.Handle(path, http.StripPrefix(path, ErrorHandler{}))
	}

	for _, mux := range muxes {
		mux.Handle(path, http.StripPrefix(path, ErrorHandler{}))
	}
}
