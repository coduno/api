package controllers

import (
	"net/http"

	"google.golang.org/appengine/datastore"

	"encoding/json"

	"github.com/coduno/app/model"
	"github.com/coduno/app/util"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

func GetProfileByKey(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if !util.CheckMethod(r, "GET") {
		return http.StatusMethodNotAllowed, nil
	}

	key, err := datastore.DecodeKey(mux.Vars(r)["key"])

	if err != nil {
		return http.StatusBadRequest, err
	}

	var profile model.Profile
	if err := datastore.Get(ctx, key, &profile); err != nil {
		return http.StatusInternalServerError, err
	}
	json.NewEncoder(w).Encode(profile.Key(key))
	return http.StatusOK, nil
}

func DeleteProfile(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if !util.CheckMethod(r, "DELETE") {
		return http.StatusMethodNotAllowed, nil
	}

	key, err := datastore.DecodeKey(mux.Vars(r)["key"])

	if err != nil {
		return http.StatusBadRequest, err
	}

	if err = datastore.Delete(ctx, key); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
