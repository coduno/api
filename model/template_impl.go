// This file was automatically generated from
//
//	template.go
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

type Templates []Template

type KeyedTemplate struct {
	*Template
	Key *datastore.Key
}

func (ƨ *Template) Key(key *datastore.Key) *KeyedTemplate {
	return &KeyedTemplate{
		Template: ƨ,
		Key:      key,
	}
}

func (ƨ Templates) Key(keys []*datastore.Key) (keyed []KeyedTemplate) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedTemplate, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedTemplate{
			Template: &ƨ[i],
			Key:      keys[i],
		}
	}
	return
}
