package controllers

import (
	"errors"
	"net/http"

	"golang.org/x/net/context"

	"github.com/coduno/api/db"
	"github.com/coduno/api/model"
)

func init() {
	router.Handle("/companies", SimpleContextHandlerFunc(PostCompany)).Methods("POST")
	router.Handle("/companies/{key}/challenges", ContextHandlerFunc(GetChallengesForCompany))
	router.Handle("/companies/{key}/users", ContextHandlerFunc(GetUsersByCompany))
}

// PostCompany creates a new company after validating by key.
func PostCompany(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var c model.Company
	if err := decode(r, &c); err != nil {
		respond(ctx, w, r, http.StatusBadRequest, err)
		return
	}
	if !c.IsValid() {
		respond(ctx, w, r, http.StatusBadRequest, db.NotValid)
		return
	}

	cs := db.NewCompanyService(ctx)
	existing, err := cs.GetByAddress(c.Address.Address)
	if err != nil {
		respond(ctx, w, r, http.StatusInternalServerError, err)
		return
	}
	if existing != nil {
		respond(ctx, w, r, http.StatusConflict, errors.New("already registered"))
		return
	}

	kc, err := cs.Save(c)
	if err != nil {
		respond(ctx, w, r, http.StatusInternalServerError, err)
		return
	}

	respond(ctx, w, r, http.StatusOK, kc)
}
