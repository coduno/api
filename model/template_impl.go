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
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
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

// Save will put this Template into Datastore using the given key.
func (ƨ Template) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Template", nil), &ƨ)
}

// SaveWithParent can be used to save this Template as child of another
// entity.
// This will error if parent == nil.
func (ƨ Template) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Template", parent), &ƨ)
}

// NewQueryForTemplate prepares a datastore.Query that can be
// used to query entities of type Template.
func NewQueryForTemplate() *datastore.Query {
	return datastore.NewQuery("Template")
}

type TemplateHandler struct{}

func (ƨ TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.URL.Path == "" {
		var results Templates
		keys, _ := NewQueryForTemplate().GetAll(ctx, &results)
		json.NewEncoder(w).Encode(results.Key(keys))
		return
	}

	k, _ := datastore.DecodeKey(r.URL.Path)
	var entity Template
	datastore.Get(ctx, k, &entity)
	json.NewEncoder(w).Encode(entity)
}

func ServeTemplate(prefix string, muxes ...*http.ServeMux) {
	path := prefix + "Template" + "/"

	if len(muxes) == 0 {
		http.Handle(path, http.StripPrefix(path, TemplateHandler{}))
	}

	for _, mux := range muxes {
		mux.Handle(path, http.StripPrefix(path, TemplateHandler{}))
	}
}
