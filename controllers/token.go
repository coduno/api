package controllers

import (
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/app/models"
	"github.com/coduno/app/util"
	"github.com/gorilla/mux"
)

// CheckToken checks the token from the client and if there exists a fingerprint
// with that token in the database send the data.
func CheckToken(w http.ResponseWriter, r *http.Request, c context.Context) (createSession bool) {
	createSession = false
	if !util.CheckMethod(w, r, "GET") {
		return
	}
	token := mux.Vars(r)["token"]
	q := datastore.NewQuery("fingerprints").Filter("Token = ", token).Limit(1)
	var fingerprints []models.Fingerprint
	_, err := q.GetAll(c, &fingerprints)
	if err != nil {
		http.Error(w, "Datastore error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(fingerprints) != 1 {
		http.Error(w, "You are unauthorized to login!", http.StatusUnauthorized)
		return
	}
	var challenge models.Challenge
	err = datastore.Get(c, fingerprints[0].Challenge, &challenge)
	if err != nil {
		http.Error(w, "Datastore error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	util.WriteEntity(w, fingerprints[0].Challenge, challenge)
	return true
}
