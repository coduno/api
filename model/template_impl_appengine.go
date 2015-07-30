// This file was automatically generated from
//
//	template.go
//
// by
//
//	generator
//
// at
//
//	2015-07-30T17:21:42+03:00
//
// Do not edit it!

package model

import (
	"errors"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// Save will put this Template into Datastore using the given key.
func (ƨ Template) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "templates", nil), &ƨ)
}

// SaveWithParent can be used to save this Template as child of another
// entity.
// This will error if parent == nil.
func (ƨ Template) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "templates", parent), &ƨ)
}

// NewQueryForTemplate prepares a datastore.Query that can be
// used to query entities of type Template.
func NewQueryForTemplate() *datastore.Query {
	return datastore.NewQuery("templates")
}

type TemplateHandler struct{}

func (ƨ TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.URL.Path == "" {
		var results Templates
		keys, _ := NewQueryForTemplate().GetAll(ctx, &results)
		results.Write(w, keys)
		return
	}

	k, _ := datastore.DecodeKey(r.URL.Path)
	var ƨ_ Template
	datastore.Get(ctx, k, &ƨ_)
	ƨ_.Write(w, k)
}

func ServeTemplate(prefix string, muxes ...*http.ServeMux) {
	path := prefix + "templates" + "/"

	if len(muxes) == 0 {
		http.Handle(path, http.StripPrefix(path, TemplateHandler{}))
	}

	for _, mux := range muxes {
		mux.Handle(path, http.StripPrefix(path, TemplateHandler{}))
	}
}
