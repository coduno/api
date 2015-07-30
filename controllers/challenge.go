package controllers

import (
	"net/http"

	"github.com/coduno/engine/model"
	"github.com/coduno/engine/passenger"
	"github.com/coduno/engine/util"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// GetChallengeByID loads a challenge by id
func GetChallengeByID(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if !util.CheckMethod(w, r, "GET") {
		return
	}

	p, ok := passenger.FromContext(ctx)

	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key, err := datastore.DecodeKey(mux.Vars(r)["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var challenge model.Challenge

	err = datastore.Get(ctx, key, &challenge)
	if err != nil {
		http.Error(w, "Datastore err"+err.Error(), http.StatusInternalServerError)
		return
	}

	if parent := p.UserKey.Parent(); parent == nil {
		// The current user is a coder so we must also create a result.
		challenge.Write(w, key)
	} else {
		// TODO(pbochis) : If a company representative user makes the request we also include Tasks in the response.
		challenge.Write(w, key)
	}
}

// GetChallengesForCompany queries all the challenges defined  by a company.
func GetChallengesForCompany(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err error

	if !util.CheckMethod(w, r, "GET") {
		return
	}

	key, err := datastore.DecodeKey(mux.Vars(r)["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := model.NewQueryForChallenge().Filter("Company=", key)

	var challenges model.Challenges

	keys, err := q.GetAll(ctx, &challenges)

	if err != nil {
		http.Error(w, "Internal Server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	values := make([]interface{}, len(challenges))
	for i := range challenges {
		values[i] = challenges[i]
	}
	challenges.Write(w, keys)
}
