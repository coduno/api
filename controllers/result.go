package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/coduno/app/model"
	"github.com/coduno/app/util"
	"github.com/coduno/app/util/passenger"
	"github.com/gorilla/mux"
	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
)

// CreateResult saves a new result when a coder starts a challenge.
func CreateResult(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if !util.CheckMethod(r, "POST") {
		return http.StatusMethodNotAllowed, nil
	}

	var body = struct {
		ChallengeKey string
	}{}

	p, ok := passenger.FromContext(ctx)

	if !ok {
		return http.StatusUnauthorized, nil
	}

	var profiles model.Profiles
	keys, err := model.NewQueryForProfile().
		Ancestor(p.UserKey).
		GetAll(ctx, &profiles)

	if len(keys) != 1 {
		return http.StatusInternalServerError, err
	}

	err = json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		return http.StatusBadRequest, err
	}

	key, err := datastore.DecodeKey(body.ChallengeKey)

	if err != nil {
		return http.StatusBadRequest, err
	}
	result := model.Result{Challenge: key}
	key, err = result.SaveWithParent(ctx, keys[0])
	if err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(result.Key(key))
	return http.StatusOK, nil
}

// GetResultsByChallenge queries the results for a certain challenge to be reviewed by a company.
func GetResultsByChallenge(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if !util.CheckMethod(r, "POST") {
		return http.StatusMethodNotAllowed, nil
	}

	key, err := datastore.DecodeKey(mux.Vars(r)["key"])

	if err != nil {
		return http.StatusBadRequest, err
	}

	var results model.Results
	keys, err := model.NewQueryForResult().
		Filter("Challenge=", key).
		GetAll(ctx, &results)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(results.Key(keys))
	return http.StatusOK, nil
}
