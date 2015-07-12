package controllers

import (
	"encoding/json"
	"net/http"

	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"

	"github.com/coduno/app/models"
	"github.com/coduno/app/util"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/m4rw3r/uuid"
)

// CheckToken checks the token from the client and if there exists a fingerprint
// with that token in the database send the data.
func CheckToken(w http.ResponseWriter, r *http.Request, c context.Context, store sessions.Store) {
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
	// TODO(victorbalan): Map the correct id
	challenge.EntityID = fingerprints[0].Challenge.StringID()

	json, err := json.Marshal(challenge)
	if err != nil {
		http.Error(w, "Json marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	initSession(w, r, store)
	w.Write([]byte(json))
}

func initSession(w http.ResponseWriter, r *http.Request, store sessions.Store) {
	id, _ := uuid.V4()
	sessionID := id.String()

	session, _ := store.Get(r, sessionID)
	session.Options.Path = "/"
	session.Options.MaxAge = 12 * 3600
	session.Save(r, w)
	w.Header().Set("Cookie", sessionID)
}
