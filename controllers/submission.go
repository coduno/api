package controllers

import (
	"encoding/json"
	"net/http"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/engine/model"
	"github.com/coduno/engine/passenger"
	"github.com/coduno/engine/util"
	"github.com/gorilla/mux"

	"golang.org/x/net/context"
)

// PostSubmission creates a new submission
func PostSubmission(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	p, ok := passenger.FromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if !util.CheckMethod(w, r, "POST") {
		return
	}
	var submission model.Submission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		http.Error(w, "Unmarshal error:"+err.Error(), http.StatusInternalServerError)
		return
	}
	resultKey, err := datastore.DecodeKey(mux.Vars(r)["id"])

	if !util.HasParent(p.UserKey, resultKey) {
		http.Error(w, "Cannot submit answer for other users", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Key decoding: "+err.Error(), http.StatusInternalServerError)
		return
	}
	key, err := submission.SaveWithParent(ctx, resultKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	submission.Write(w, key)
}
