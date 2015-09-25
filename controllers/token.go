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

func init() {
	router.Handle("/tokens", ContextHandlerFunc(Tokens))
	router.Handle("/tokens/collect", ContextHandlerFunc(CollectTokens))
}

// Tokens will create new Tokens for the user.
func Tokens(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, nil
	}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	var body model.Token
	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return http.StatusBadRequest, err
	}

	value, err := p.IssueToken(ctx, &body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

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

	if err = json.NewEncoder(w).Encode(result); err != nil {
		return http.StatusInternalServerError, err
	}
	return
}

// CollectTokens runs a query against Datastore to find expired
// Tokens and deletes them.
func CollectTokens(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	keys, err := model.NewQueryForToken().
		Filter("Expiry <", time.Now()).
		KeysOnly().
		GetAll(ctx, nil)

	if err != nil {
		log.Warningf(ctx, "garbage collecting tokens failed: %s", err.Error())
		return http.StatusInternalServerError, err
	}

	err = datastore.DeleteMulti(ctx, keys)

	if err != nil {
		log.Warningf(ctx, "garbage collecting tokens failed: %s", err.Error())
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
