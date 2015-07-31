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
	"google.golang.org/appengine/datastore"
)

type AccessTokens []AccessToken

type KeyedAccessToken struct{
	AccessToken *AccessToken
	Key      *datastore.Key
}

func (ƨ *AccessToken) Key(key *datastore.Key) (*KeyedAccessToken) {
	return &KeyedAccessToken{
		AccessToken: ƨ,
		Key:      key,
	}
}

func (ƨ AccessTokens) Key(keys []*datastore.Key) (keyed []KeyedAccessToken) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedAccessToken, 0, len(ƨ))
	for i := range keyed {
		keyed = append(keyed, KeyedAccessToken{
			AccessToken: &ƨ[i],
			Key:      keys[i],
		})
	}
	return
}
