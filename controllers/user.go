package controllers

import (
	"net/http"

	"encoding/json"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"github.com/coduno/api/util/passenger"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// GetUsersByCompany queries the user accounts belonging to a company.
func GetUsersByCompany(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if !util.CheckMethod(r, "GET") {
		return http.StatusMethodNotAllowed, nil
	}
	key, err := datastore.DecodeKey(r.URL.Query()["result"][0])
	if err != nil {
		return http.StatusBadRequest, err
	}
	var users model.Users
	keys, err := model.NewQueryForUser().
		Ancestor(key).
		GetAll(ctx, &users)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	json.NewEncoder(w).Encode(users.Key(keys))
	return http.StatusOK, nil
}

func GetCompanyByUser(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if !util.CheckMethod(r, "GET") {
		return http.StatusMethodNotAllowed, nil
	}
	p, ok := passenger.FromContext(ctx)

	if !ok {
		return http.StatusUnauthorized, nil
	}
	key := p.UserKey.Parent()
	if key == nil {
		return http.StatusUnauthorized, nil
	}
	// The account is associated with a company, so we return it.
	var company model.Company
	if err := datastore.Get(ctx, key, &company); err != nil {
		return http.StatusInternalServerError, err
	}
	json.NewEncoder(w).Encode(company.Key(key))
	return http.StatusOK, nil
}

func WhoAmI(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if !util.CheckMethod(r, "GET") {
		return http.StatusMethodNotAllowed, nil
	}
	p, ok := passenger.FromContext(ctx)

	if !ok {
		return http.StatusUnauthorized, nil
	}
	json.NewEncoder(w).Encode(p.User.Key(p.UserKey))
	return http.StatusOK, nil
}
