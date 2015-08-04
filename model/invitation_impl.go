// This file was automatically generated from
//
//	invitation.go
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

type Invitations []Invitation

type KeyedInvitation struct {
	*Invitation
	Key *datastore.Key
}

func (ƨ *Invitation) Key(key *datastore.Key) *KeyedInvitation {
	return &KeyedInvitation{
		Invitation: ƨ,
		Key:        key,
	}
}

func (ƨ Invitations) Key(keys []*datastore.Key) (keyed []KeyedInvitation) {
	if len(keys) != len(ƨ) {
		panic("Key() called on an slice with len(keys) != len(slice)")
	}

	keyed = make([]KeyedInvitation, len(ƨ))
	for i := range keyed {
		keyed[i] = KeyedInvitation{
			Invitation: &ƨ[i],
			Key:        keys[i],
		}
	}
	return
}

// Save will put this Invitation into Datastore using the given key.
func (ƨ Invitation) Save(ctx context.Context, key ...*datastore.Key) (*datastore.Key, error) {
	if len(key) > 1 {
		panic("zero or one key expected")
	}

	if len(key) == 1 && key[0] != nil {
		return datastore.Put(ctx, key[0], &ƨ)
	}

	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Invitation", nil), &ƨ)
}

// SaveWithParent can be used to save this Invitation as child of another
// entity.
// This will error if parent == nil.
func (ƨ Invitation) SaveWithParent(ctx context.Context, parent *datastore.Key) (*datastore.Key, error) {
	if parent == nil {
		return nil, errors.New("parent key is nil, expected a valid key")
	}
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Invitation", parent), &ƨ)
}

// NewQueryForInvitation prepares a datastore.Query that can be
// used to query entities of type Invitation.
func NewQueryForInvitation() *datastore.Query {
	return datastore.NewQuery("Invitation")
}

type InvitationHandler struct{}

func (ƨ InvitationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.URL.Path == "" {
		var results Invitations
		keys, _ := NewQueryForInvitation().GetAll(ctx, &results)
		json.NewEncoder(w).Encode(results.Key(keys))
		return
	}

	k, _ := datastore.DecodeKey(r.URL.Path)
	var entity Invitation
	datastore.Get(ctx, k, &entity)
	json.NewEncoder(w).Encode(entity)
}

func ServeInvitation(prefix string, muxes ...*http.ServeMux) {
	path := prefix + "Invitation" + "/"

	if len(muxes) == 0 {
		http.Handle(path, http.StripPrefix(path, InvitationHandler{}))
	}

	for _, mux := range muxes {
		mux.Handle(path, http.StripPrefix(path, InvitationHandler{}))
	}
}
