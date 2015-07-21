package controllers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/app/models"
	"github.com/coduno/app/util"
	"github.com/gorilla/mux"
)

// LoadFingerprintsByCompanyID -
func LoadFingerprintsByCompanyID(w http.ResponseWriter, r *http.Request, c context.Context) {
	if !util.CheckMethod(w, r, "GET") {
		return
	}

	companyID := mux.Vars(r)["companyId"]
	key, _ := datastore.DecodeKey(companyID)
	q := datastore.NewQuery("challenges").Filter("Company = ", key).KeysOnly()
	keys, err := q.GetAll(c, nil)
	if err != nil {
		http.Error(w, "Datastore error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	q = datastore.NewQuery("fingerprints")
	for _, val := range keys {
		q.Filter("Challenge = ", val)
	}
	var fingerprints []models.Fingerprint
	keys, err = q.GetAll(c, &fingerprints)
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
