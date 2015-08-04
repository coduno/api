// This file was automatically generated from
//
//	code_submission.go
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

type CodeSubmissions []CodeSubmission

type KeyedCodeSubmission struct {
	*CodeSubmission
	Key *datastore.Key
}

func (ƨ *CodeSubmission) Key(key *datastore.Key) *KeyedCodeSubmission {
	return &KeyedCodeSubmission{
		CodeSubmission: ƨ,
		Key:            key,
	}
}

func (ƨ CodeSubmissions) Key(keys []*datastore.Key) (keyed []KeyedCodeSubmission) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedCodeSubmission, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedCodeSubmission{
			CodeSubmission: &ƨ[i],
			Key:            keys[i],
		}
	}
	return
}

// Save will put this CodeSubmission into Datastore using the given key.
func (ƨ CodeSubmission) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "CodeSubmission", nil), &ƨ)
}

// SaveWithParent can be used to save this CodeSubmission as child of another
// entity.
// This will error if parent == nil.
func (ƨ CodeSubmission) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "CodeSubmission", parent), &ƨ)
}

// NewQueryForCodeSubmission prepares a datastore.Query that can be
// used to query entities of type CodeSubmission.
func NewQueryForCodeSubmission() *datastore.Query {
	return datastore.NewQuery("CodeSubmission")
}

type CodeSubmissionHandler struct{}

func (ƨ CodeSubmissionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.URL.Path == "" {
		var results CodeSubmissions
		keys, _ := NewQueryForCodeSubmission().GetAll(ctx, &results)
		json.NewEncoder(w).Encode(results.Key(keys))
		return
	}

	k, _ := datastore.DecodeKey(r.URL.Path)
	var entity CodeSubmission
	datastore.Get(ctx, k, &entity)
	json.NewEncoder(w).Encode(entity)
}

func ServeCodeSubmission(prefix string, muxes ...*http.ServeMux) {
	path := prefix + "CodeSubmission" + "/"

	if len(muxes) == 0 {
		http.Handle(path, http.StripPrefix(path, CodeSubmissionHandler{}))
	}

	for _, mux := range muxes {
		mux.Handle(path, http.StripPrefix(path, CodeSubmissionHandler{}))
	}
}
