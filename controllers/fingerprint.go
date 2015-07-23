package controllers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/app/models"
	"github.com/coduno/app/util"
	"github.com/gorilla/mux"
	"github.com/m4rw3r/uuid"
)

//CreateFingerprint from the request body
func CreateFingerprint(w http.ResponseWriter, r *http.Request, ctx context.Context) {

	var err error

	if !util.CheckMethod(w, r, "POST") {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	firstName := mux.Vars(r)["firstName"]
	lastName := mux.Vars(r)["lastName"]
	email := mux.Vars(r)["email"]

	coder := models.Coder{FirstName: firstName, LastName: lastName, Email: email}

	coderKey, err := coder.Save(ctx)

	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	challangeID := mux.Vars(r)["challangeId"]
	challangeKey := datastore.NewKey(ctx, models.ChallangeKind, challangeID, 0, nil)

	token, err := uuid.V4()

	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fingerprint := models.Fingerprint{Coder: coderKey, Challenge: challangeKey, Token: token.String()}
	_, err = fingerprint.Save(ctx)

	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
	}

	response, err := json.Marshal(fingerprint)

	if err != nil {
		http.Error(w, "Json marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO(pbochis): This is where we will send an e-mail to the candidate with
	// somthing like "cod.uno/fingerprint/:token".

	w.Write([]byte(response))
	return
}

// LoadFingerprintsByCompanyID -
func LoadFingerprintsByCompanyID(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if !util.CheckMethod(w, r, "GET") {
		return
	}

	companyID := mux.Vars(r)["companyId"]
	key, _ := datastore.DecodeKey(companyID)
	q := datastore.NewQuery("challenges").Filter("Company = ", key).KeysOnly()
	keys, err := q.GetAll(ctx, nil)
	if err != nil {
		http.Error(w, "Datastore error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	q = datastore.NewQuery("fingerprints")
	for _, val := range keys {
		q.Filter("Challenge = ", val)
	}
	var fingerprints []models.Fingerprint
	keys, err = q.GetAll(ctx, &fingerprints)
	if err != nil {
		http.Error(w, "Datastore error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	for i := 0; i < len(keys); i++ {
		fingerprints[i].EntityID = keys[i].Encode()
	}
	json, err := json.Marshal(fingerprints)
	if err != nil {
		http.Error(w, "Json marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(json))
}
