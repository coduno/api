package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"github.com/coduno/api/util/passenger"
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
		return http.StatusInternalServerError, errors.New("Profile not found")
	}

	err = json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		return http.StatusBadRequest, err
	}

	key, err := datastore.DecodeKey(body.ChallengeKey)

	if err != nil {
		return http.StatusBadRequest, err
	}

	var results []model.Result

	resultKeys, err := model.NewQueryForResult().
		Ancestor(keys[0]).
		Filter("Challenge = ", key).
		Limit(1).
		GetAll(ctx, &results)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if len(resultKeys) == 1 {
		json.NewEncoder(w).Encode(results[0].Key(resultKeys[0]))
		return http.StatusOK, nil
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
