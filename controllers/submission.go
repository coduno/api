package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/app/model"
	"github.com/coduno/app/util"
	"github.com/coduno/engine/passenger"
	"github.com/gorilla/mux"

	"golang.org/x/net/context"
)

// PostSubmission creates a new submission
func PostSubmission(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, errors.New("Unauthorized request")
	}
	if err = util.CheckMethod(r, "GET"); err != nil {
		return http.StatusMethodNotAllowed, err
	}
	var submission model.Submission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		return http.StatusInternalServerError, err
	}
	resultKey, err := datastore.DecodeKey(mux.Vars(r)["id"])

	if !util.HasParent(p.UserKey, resultKey) {
		return http.StatusBadRequest, errors.New("Cannot submit answer for other users")
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	key, err := submission.SaveWithParent(ctx, resultKey)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	submission.Write(w, key)
	return http.StatusCreated, nil
}
