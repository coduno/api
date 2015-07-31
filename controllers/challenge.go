package controllers

import (
	"net/http"

	"github.com/coduno/app/model"
	"github.com/coduno/engine/passenger"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// ChallengeByKey loads a challenge by key.
func ChallengeByKey(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	key, err := datastore.DecodeKey(mux.Vars(r)["key"])
	if err != nil {
		return http.StatusInternalServerError, err
	}

	var challenge model.Challenge

	err = datastore.Get(ctx, key, &challenge)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if parent := p.UserKey.Parent(); parent == nil {
		// The current user is a coder so we must also create a result.
		challenge.Write(w, key)
	} else {
		// TODO(pbochis): If a company representativemakes the request
		// we also include Tasks in the response.
		challenge.Write(w, key)
	}
	return http.StatusOK, nil
}

// GetChallengesForCompany queries all the challenges defined  by a company.
func GetChallengesForCompany(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	_, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	key, err := datastore.DecodeKey(mux.Vars(r)["key"])
	if err != nil {
		return http.StatusInternalServerError, err
	}

	var challenges model.Challenges

	keys, err := model.NewQueryForChallenge().
		Ancestor(key).
		GetAll(ctx, &challenges)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	challenges.Write(w, keys)
	return http.StatusOK, nil
}
