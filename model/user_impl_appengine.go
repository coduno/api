// This file was automatically generated from
//
//	user.go
//
// by
//
//	generator -c Profile
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

// Save will put this User into Datastore using the given key.
func (ƨ User) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "users", nil), &ƨ)
}

// SaveWithParent can be used to save this User as child of another
// entity.
// This will error if parent == nil.
func (ƨ User) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "users", parent), &ƨ)
}

// NewQueryForUser prepares a datastore.Query that can be
// used to query entities of type User.
func NewQueryForUser() *datastore.Query {
	return datastore.NewQuery("users")
}

type UserHandler struct{}

func (ƨ UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.URL.Path == "" {
		var results Users
		keys, _ := NewQueryForUser().GetAll(ctx, &results)
		json.NewEncoder(w).Encode(results.Key(keys))
		return
	}

	k, _ := datastore.DecodeKey(r.URL.Path)
	var entity User
	datastore.Get(ctx, k, &entity)
	json.NewEncoder(w).Encode(entity)
}

func ServeUser(prefix string, muxes ...*http.ServeMux) {
	path := prefix + "users" + "/"

	if len(muxes) == 0 {
		http.Handle(path, http.StripPrefix(path, UserHandler{}))
	}

	for _, mux := range muxes {
		mux.Handle(path, http.StripPrefix(path, UserHandler{}))
	}
}
