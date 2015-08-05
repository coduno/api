// This file was automatically generated from
//
//	company.go
//
// by
//
//	generator -c Challenge,User
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

// CompanyKind is the kind used in Datastore to store entities of type Company.
const CompanyKind = "Company"

// Companys is just a slice of Company.
type Companys []Company

// KeyedCompany is a struct that embeds Company and also contains a Key, mainly used for encoding to JSON.
type KeyedCompany struct {
	*Company
	Key *datastore.Key
}

// Key is a shorthand to fill a KeyedCompany with an entity and it's key.
func (ƨ *Company) Key(key *datastore.Key) *KeyedCompany {
	return &KeyedCompany{
		Company: ƨ,
		Key:     key,
	}
}

// Key is a shorthand to fill a slice of KeyedCompany with some entities alongside their keys.
func (ƨ Companys) Key(keys []*datastore.Key) (keyed []KeyedCompany) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedCompany, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedCompany{
			Company: &ƨ[i],
			Key:     keys[i],
		}
	}
	return
}

// Save will put this Company into Datastore using the given key.
func (ƨ Company) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Company", nil), &ƨ)
}

// SaveWithParent can be used to save this Company as child of another
// entity.
// This will error if parent == nil.
func (ƨ Company) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Company", parent), &ƨ)
}

// NewQueryForCompany prepares a datastore.Query that can be
// used to query entities of type Company.
func NewQueryForCompany() *datastore.Query {
	return datastore.NewQuery("Company")
}
