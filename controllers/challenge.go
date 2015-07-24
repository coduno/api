package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/coduno/app/models"
	"github.com/coduno/app/util"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// GetChallangesForCompany -
func GetChallangesForCompany(w http.ResponseWriter, r *http.Request, ctx context.Context) {

	var err error

	if !util.CheckMethod(w, r, "GET") {
		return
	}

	companyKey := mux.Vars(r)["companyId"]

	if companyKey == "" {
		http.Error(w, "You need to provide a company id", http.StatusInternalServerError)
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

	for i := 0; i < len(keys); i++ {
		challenges[i].EntityID = keys[i].Encode()
	}

	response := make(map[string]interface{})
	response["challanges"] = challenges
	body, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Internal Server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(body)
	return
}
