package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/coduno/app/model"
	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
)

func CreateResult(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	var body = struct {
		ChallengeID string
	}{}

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return http.StatusInternalServerError, err
	}

	key, err := datastore.DecodeKey(body.ChallengeID)
	if err != nil {
		return http.StatusBadRequest, err
	}

	result := model.Result{Challenge: key}
	key, err = result.Save(ctx)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	result.Write(w, key)
	return http.StatusOK, nil
}
