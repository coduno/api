// This file was automatically generated from
//
//	profile.go
//
// by
//
//	generator -c Result
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

// Save will put this Profile into Datastore using the given key.
func (ƨ Profile) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "profiles", nil), &ƨ)
}

// SaveWithParent can be used to save this Profile as child of another
// entity.
// This will error if parent == nil.
func (ƨ Profile) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "profiles", parent), &ƨ)
}

// NewQueryForProfile prepares a datastore.Query that can be
// used to query entities of type Profile.
func NewQueryForProfile() *datastore.Query {
	return datastore.NewQuery("profiles")
}

type ProfileHandler struct{}

func (ƨ ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.URL.Path == "" {
		var results Profiles
		keys, _ := NewQueryForProfile().GetAll(ctx, &results)
		results.Write(w, keys)
		return
	}

	k, _ := datastore.DecodeKey(r.URL.Path)
	var ƨ_ Profile
	datastore.Get(ctx, k, &ƨ_)
	ƨ_.Write(w, k)
}

func ServeProfile(prefix string, muxes ...*http.ServeMux) {
	path := prefix + "profiles" + "/"

	if len(muxes) == 0 {
		http.Handle(path, http.StripPrefix(path, ProfileHandler{}))
	}

	for _, mux := range muxes {
		mux.Handle(path, http.StripPrefix(path, ProfileHandler{}))
	}
}
