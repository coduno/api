package controllers

import (
	"net/http"

	"github.com/coduno/app/models"
	"github.com/coduno/app/util"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// GetChallengeByID loads a challenge by id
func GetChallengeByID(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if !util.CheckMethod(w, r, "GET") {
		return
	}

	key, err := datastore.DecodeKey(mux.Vars(r)["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var challenge models.Challenge
	err = datastore.Get(ctx, key, &challenge)
	if err != nil {
		http.Error(w, "Datastore err"+err.Error(), http.StatusInternalServerError)
		return
	}
	util.WriteEntity(w, key, challenge)
}

func GetChallengesForCompany(w http.ResponseWriter, r *http.Request, ctx context.Context) {

	var err error

	if !util.CheckMethod(w, r, "GET") {
		return
	}

	companyKey := r.URL.Query()["company"][0]

	if companyKey == "" {
		http.Error(w, "missing parameter 'company'", http.StatusInternalServerError)
		return
	}

	key, err := datastore.DecodeKey(companyKey)

	if err != nil {
		http.Error(w, "Invalid company", http.StatusInternalServerError)
		return
	}

	q := datastore.NewQuery(models.ChallengeKind).Filter("Company=", key)

	var challenges []models.Challenge

	keys, err := q.GetAll(ctx, &challenges)

	if err != nil {
		http.Error(w, "Internal Server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	values := make([]interface{}, len(challenges))
	for i := range challenges {
		values[i] = challenges[i]
	}
	util.WriteEntities(w, keys, values)
}
