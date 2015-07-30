package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/coduno/engine/model"
	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
)

type ResultData struct {
	ChallengeId string
}

func CreateResult(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var resultData ResultData

	err := decoder.Decode(&resultData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key, err := datastore.DecodeKey(resultData.ChallengeId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := model.Result{Challenge: key}

	key, err = result.Save(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result.Write(w, key)
}
