package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util/passenger"

	"golang.org/x/net/context"
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
