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
	"google.golang.org/appengine/datastore"
)

type Companys []Company

type KeyedCompany struct {
	*Company
	Key *datastore.Key
}

func (ƨ *Company) Key(key *datastore.Key) *KeyedCompany {
	return &KeyedCompany{
		Company: ƨ,
		Key:     key,
	}
}

func (ƨ Companys) Key(keys []*datastore.Key) (keyed []KeyedCompany) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedCompany, 0, len(ƨ))
	for i := range keyed {
		keyed = append(keyed, KeyedCompany{
			Company: &ƨ[i],
			Key:     keys[i],
		})
	}
	return
}
