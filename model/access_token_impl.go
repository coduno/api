// This file was automatically generated from
//
//	access_token.go
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

const AccessTokenKind = "AccessToken"

type AccessTokens []AccessToken

type KeyedAccessToken struct {
	*AccessToken
	Key *datastore.Key
}

func (ƨ *AccessToken) Key(key *datastore.Key) *KeyedAccessToken {
	return &KeyedAccessToken{
		AccessToken: ƨ,
		Key:         key,
	}
}

func (ƨ AccessTokens) Key(keys []*datastore.Key) (keyed []KeyedAccessToken) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedAccessToken, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedAccessToken{
			AccessToken: &ƨ[i],
			Key:         keys[i],
		}
	}
	return
}

// Save will put this AccessToken into Datastore using the given key.
func (ƨ AccessToken) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "AccessToken", nil), &ƨ)
}

// SaveWithParent can be used to save this AccessToken as child of another
// entity.
// This will error if parent == nil.
func (ƨ AccessToken) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "AccessToken", parent), &ƨ)
}

// NewQueryForAccessToken prepares a datastore.Query that can be
// used to query entities of type AccessToken.
func NewQueryForAccessToken() *datastore.Query {
	return datastore.NewQuery("AccessToken")
}

type AccessTokenHandler struct{}

func (ƨ AccessTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.URL.Path == "" {
		var results AccessTokens
		keys, _ := NewQueryForAccessToken().GetAll(ctx, &results)
		json.NewEncoder(w).Encode(results.Key(keys))
		return
	}

	k, _ := datastore.DecodeKey(r.URL.Path)
	var entity AccessToken
	datastore.Get(ctx, k, &entity)
	json.NewEncoder(w).Encode(entity)
}

func ServeAccessToken(prefix string, muxes ...*http.ServeMux) {
	path := prefix + "AccessToken" + "/"

	if len(muxes) == 0 {
		http.Handle(path, http.StripPrefix(path, AccessTokenHandler{}))
	}

	for _, mux := range muxes {
		mux.Handle(path, http.StripPrefix(path, AccessTokenHandler{}))
	}
}
