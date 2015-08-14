// This file was automatically generated from
//
//	token.go
//
// by
//
//	generator
//
// DO NOT EDIT

package model

import (
	"errors"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// TokenKind is the kind used in Datastore to store entities Token entities.
const TokenKind = "Token"

// Tokens is just a slice of Token.
type Tokens []Token

// KeyedToken is a struct that embeds Token and also contains a Key, mainly used for encoding to JSON.
type KeyedToken struct {
	*Token
	Key *datastore.Key
}

// Key is a shorthand to fill a KeyedToken with an entity and it's key.
func (ƨ *Token) Key(key *datastore.Key) *KeyedToken {
	return &KeyedToken{
		Token: ƨ,
		Key:   key,
	}
}

// Key is a shorthand to fill a slice of KeyedToken with some entities alongside their keys.
func (ƨ Tokens) Key(keys []*datastore.Key) (keyed []KeyedToken) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedToken, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedToken{
			Token: &ƨ[i],
			Key:   keys[i],
		}
	}
	return
}

// Save will put this Token into Datastore using the given key.
func (ƨ Token) Save(ctx context.Context, key *datastore.Key) (*datastore.Key, error) {
	if key != nil {
		return datastore.Put(ctx, key, &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Token", nil), &ƨ)
}

// SaveWithParent can be used to save this Token as child of another
// entity.
// This will error if parent == nil.
func (ƨ Token) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Token", parent), &ƨ)
}

// NewQueryForToken prepares a datastore.Query that can be
// used to query entities of type Token.
func NewQueryForToken() *datastore.Query {
	return datastore.NewQuery("Token")
}
