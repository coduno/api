package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util/passenger"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// AccessTokens will create new AccessTokens for the user.
func AccessTokens(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, nil
	}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	var body model.AccessToken

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return http.StatusBadRequest, err
	}

	value, err := p.IssueToken(ctx, &body)

	var result = struct {
		Value       string
		Creation    time.Time
		Expiry      time.Time
		Description string
	}{
		Value:       value,
		Creation:    body.Creation,
		Expiry:      body.Expiry,
		Description: body.Description,
	}

	json.NewEncoder(w).Encode(result)
	return
}

// CollectAccessTokens runs a query against Datastore to find expired
// AccessTokens and deletes them.
func CollectAccessTokens(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	keys, err := model.NewQueryForAccessToken().
		Filter("Expiry <", time.Now()).
		KeysOnly().
		GetAll(ctx, nil)

	if err != nil {
		log.Warningf(ctx, "garbage collecting access tokens failed: %s", err.Error())
		return http.StatusInternalServerError, err
	}

	err = datastore.DeleteMulti(ctx, keys)

	if err != nil {
		log.Warningf(ctx, "garbage collecting access tokens failed: %s", err.Error())
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
