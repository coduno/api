// This file was automatically generated from
//
//	challenge.go
//
// by
//
//	generator
//
// DO NOT EDIT

package model

import (
	"errors"
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// Save will put this Challenge into Datastore using the given key.
func (ƨ Challenge) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "challenges", nil), &ƨ)
}

// SaveWithParent can be used to save this Challenge as child of another
// entity.
// This will error if parent == nil.
func (ƨ Challenge) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "challenges", parent), &ƨ)
}

// NewQueryForChallenge prepares a datastore.Query that can be
// used to query entities of type Challenge.
func NewQueryForChallenge() *datastore.Query {
	return datastore.NewQuery("challenges")
}

type ChallengeHandler struct{}

func (ƨ ChallengeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.URL.Path == "" {
		var results Challenges
		keys, _ := NewQueryForChallenge().GetAll(ctx, &results)
		json.NewEncoder(w).Encode(results.Key(keys))
		return
	}

	k, _ := datastore.DecodeKey(r.URL.Path)
	var entity Challenge
	datastore.Get(ctx, k, &entity)
	json.NewEncoder(w).Encode(entity)
}

func ServeChallenge(prefix string, muxes ...*http.ServeMux) {
	path := prefix + "challenges" + "/"

	if len(muxes) == 0 {
		http.Handle(path, http.StripPrefix(path, ChallengeHandler{}))
	}

	for _, mux := range muxes {
		mux.Handle(path, http.StripPrefix(path, ChallengeHandler{}))
	}
}
