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

// AccessTokenKind is the kind used in Datastore to store entities of type AccessToken.
const AccessTokenKind = "AccessToken"

// AccessTokens is just a slice of AccessToken.
type AccessTokens []AccessToken

// KeyedAccessToken is a struct that embeds AccessToken and also contains a Key, mainly used for encoding to JSON.
type KeyedAccessToken struct {
	*AccessToken
	Key *datastore.Key
}

// Key is a shorthand to fill a KeyedAccessToken with an entity and it's key.
func (ƨ *AccessToken) Key(key *datastore.Key) *KeyedAccessToken {
	return &KeyedAccessToken{
		AccessToken: ƨ,
		Key:         key,
	}
}

// Key is a shorthand to fill a slice of KeyedAccessToken with some entities alongside their keys.
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
