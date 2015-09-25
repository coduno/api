package controllers

import (
	"net/http"
	"strconv"
	"time"

	"google.golang.org/appengine/datastore"

	"encoding/json"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util/passenger"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

func init() {
	router.HandleFunc("/profiles/{key}", setup(GetProfileByKey))
	router.HandleFunc("/profiles/{key}", setup(DeleteProfile))
	router.HandleFunc("/profiles/{key}/challenges", setup(GetChallengesForProfile))
}

func GetProfileByKey(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "GET" {
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

func GetChallengesForProfile(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	_, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	var profileKey *datastore.Key
	if profileKey, err = datastore.DecodeKey(mux.Vars(r)["key"]); err != nil {
		return http.StatusInternalServerError, err
	}

	q := model.NewQueryForResult().
		Ancestor(profileKey)
	if finished := r.URL.Query()["finished"]; len(finished) > 0 && finished[0] == "true" {
		q = q.Filter("Finished >", time.Time{})
	}
	if order := r.URL.Query()["order"]; len(order) > 0 && order[0] != "" {
		q = q.Order(order[0])
	}

	if limitQuery := r.URL.Query()["limit"]; len(limitQuery) > 0 {
		if limit, err := strconv.Atoi(limitQuery[0]); err != nil {
			return http.StatusInternalServerError, err
		} else {
			q = q.Limit(limit)
		}
	}

	var results model.Results
	if _, err = q.GetAll(ctx, &results); err != nil {
		return http.StatusInternalServerError, err
	}

	challengeKeys := make([]*datastore.Key, len(results))
	for i, val := range results {
		challengeKeys[i] = val.Challenge
	}

	challenges := make(model.Challenges, len(challengeKeys))
	if err = datastore.GetMulti(ctx, challengeKeys, challenges); err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(challenges.Key(challengeKeys))
	return
}

func GetProfileForUser(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	_, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	var userKey *datastore.Key
	if userKey, err = datastore.DecodeKey(mux.Vars(r)["key"]); err != nil {
		return http.StatusInternalServerError, err
	}

	var profiles []model.Profile
	keys, err := model.NewQueryForProfile().
		Ancestor(userKey).
		Limit(1).
		GetAll(ctx, &profiles)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if len(keys) != 1 {
		return http.StatusNotFound, nil
	}
	json.NewEncoder(w).Encode(profiles[0].Key(keys[0]))
	return
}

func DeleteProfile(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "DELETE" {
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
