// This file was automatically generated from
//
//	user.go
//
// by
//
//	generator -c Profile
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
		results.Write(w, keys)
		return
	}

	k, _ := datastore.DecodeKey(r.URL.Path)
	var ƨ_ User
	datastore.Get(ctx, k, &ƨ_)
	ƨ_.Write(w, k)
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
