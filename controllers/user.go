package controllers

import (
	"net/http"

	"encoding/json"

	"github.com/coduno/app/model"
	"github.com/coduno/app/util"
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
