package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/coduno/api/logic"
	"github.com/coduno/api/model"
	"github.com/coduno/api/util/passenger"
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
		return http.StatusBadRequest, err
	}

	var challenge model.Challenge

	err = datastore.Get(ctx, key, &challenge)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	var u model.User
	if err := datastore.Get(ctx, p.User, &u); err != nil {
		return http.StatusInternalServerError, err
	}

	e := json.NewEncoder(w)
	if u.Company == nil {
		// The current user is a coder so we must also create a result.
		e.Encode(challenge.Key(key))
	} else {
		// TODO(pbochis): If a company representativemakes the request
		// we also include Tasks in the response.
		e.Encode(challenge.Key(key))
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

	json.NewEncoder(w).Encode(challenges.Key(keys))
	return http.StatusOK, nil
}

// CreateChallenge will put a new entity of kind Challenge to Datastore.
func CreateChallenge(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, nil
	}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	var u model.User
	if err := datastore.Get(ctx, p.User, &u); err != nil {
		return http.StatusInternalServerError, err
	}

	if u.Company == nil {
		return http.StatusUnauthorized, nil
	}

	var body = struct {
		model.Assignment
		Tasks []string
	}{}

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return http.StatusBadRequest, err
	}

	keys := make([]*datastore.Key, len(body.Tasks))
	for i := range body.Tasks {
		key, err := datastore.DecodeKey(body.Tasks[i])
		if err != nil {
			return http.StatusInternalServerError, err
		}
		keys[i] = key
	}

	challenge := model.Challenge{
		Assignment: body.Assignment,
		Resulter:   int64(logic.Average),
		Tasks:      keys,
	}

	key, err := challenge.PutWithParent(ctx, u.Company)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(challenge.Key(key))
	return http.StatusOK, nil
}
