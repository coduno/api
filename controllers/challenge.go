package controllers

import (
	"errors"
	"net/http"

	"github.com/coduno/engine/model"
	"github.com/coduno/engine/passenger"
	"github.com/coduno/engine/util"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// GetChallengeByID loads a challenge by id
func GetChallengeByID(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if err = util.CheckMethod(r, "GET"); err != nil {
		return http.StatusMethodNotAllowed, err
	}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, errors.New("Unauthorized request")
	}

	key, err := datastore.DecodeKey(mux.Vars(r)["id"])

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
		// TODO(pbochis) : If a company representative user makes the request we also include Tasks in the response.
		challenge.Write(w, key)
	}
	return http.StatusOK, nil
}

// GetChallengesForCompany queries all the challenges defined  by a company.
func GetChallengesForCompany(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if err = util.CheckMethod(r, "GET"); err != nil {
		return http.StatusMethodNotAllowed, err
	}
	_, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, errors.New("Unauthorized request")
	}
	key, err := datastore.DecodeKey(mux.Vars(r)["id"])

	if err != nil {
		return http.StatusInternalServerError, err
	}

	q := model.NewQueryForChallenge().Ancestor(key)

	var challenges model.Challenges

	keys, err := q.GetAll(ctx, &challenges)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	values := make([]interface{}, len(challenges))
	for i := range challenges {
		values[i] = challenges[i]
	}
	challenges.Write(w, keys)
	return http.StatusOK, nil
}
