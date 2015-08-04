// This file was automatically generated from
//
//	result.go
//
// by
//
//	generator -c Submission
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

const ResultKind = "Result"

type Results []Result

type KeyedResult struct {
	*Result
	Key *datastore.Key
}

func (ƨ *Result) Key(key *datastore.Key) *KeyedResult {
	return &KeyedResult{
		Result: ƨ,
		Key:    key,
	}
}

func (ƨ Results) Key(keys []*datastore.Key) (keyed []KeyedResult) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedResult, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedResult{
			Result: &ƨ[i],
			Key:    keys[i],
		}
	}
	return
}

// Save will put this Result into Datastore using the given key.
func (ƨ Result) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Result", nil), &ƨ)
}

// SaveWithParent can be used to save this Result as child of another
// entity.
// This will error if parent == nil.
func (ƨ Result) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Result", parent), &ƨ)
}

// NewQueryForResult prepares a datastore.Query that can be
// used to query entities of type Result.
func NewQueryForResult() *datastore.Query {
	return datastore.NewQuery("Result")
}

type ResultHandler struct{}

func (ƨ ResultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.URL.Path == "" {
		var results Results
		keys, _ := NewQueryForResult().GetAll(ctx, &results)
		json.NewEncoder(w).Encode(results.Key(keys))
		return
	}

	k, _ := datastore.DecodeKey(r.URL.Path)
	var entity Result
	datastore.Get(ctx, k, &entity)
	json.NewEncoder(w).Encode(entity)
}

func ServeResult(prefix string, muxes ...*http.ServeMux) {
	path := prefix + "Result" + "/"

	if len(muxes) == 0 {
		http.Handle(path, http.StripPrefix(path, ResultHandler{}))
	}

	for _, mux := range muxes {
		mux.Handle(path, http.StripPrefix(path, ResultHandler{}))
	}
}
