// This file was automatically generated from
//
//	invitation.go
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

// Save will put this Invitation into Datastore using the given key.
func (ƨ Invitation) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "invitations", nil), &ƨ)
}

// SaveWithParent can be used to save this Invitation as child of another
// entity.
// This will error if parent == nil.
func (ƨ Invitation) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "invitations", parent), &ƨ)
}

// NewQueryForInvitation prepares a datastore.Query that can be
// used to query entities of type Invitation.
func NewQueryForInvitation() *datastore.Query {
	return datastore.NewQuery("invitations")
}

type InvitationHandler struct{}

func (ƨ InvitationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.URL.Path == "" {
		var results Invitations
		keys, _ := NewQueryForInvitation().GetAll(ctx, &results)
		results.Write(w, keys)
		return
	}

	k, _ := datastore.DecodeKey(r.URL.Path)
	var ƨ_ Invitation
	datastore.Get(ctx, k, &ƨ_)
	ƨ_.Write(w, k)
}

func ServeInvitation(prefix string, muxes ...*http.ServeMux) {
	path := prefix + "invitations" + "/"

	if len(muxes) == 0 {
		http.Handle(path, http.StripPrefix(path, InvitationHandler{}))
	}

	for _, mux := range muxes {
		mux.Handle(path, http.StripPrefix(path, InvitationHandler{}))
	}
}
