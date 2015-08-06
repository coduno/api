package controllers

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"google.golang.org/appengine/memcache"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util/passenger"
	"github.com/coduno/api/util/password"

	"golang.org/x/net/context"
)

const maxValidity = time.Hour * 24 * 30 * 2

// AccessTokens will create new AccessTokens for the user.
func AccessTokens(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, nil
	}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	var body = struct {
		Scopes      []string
		Description string
		Expiry      time.Time
	}{}

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return http.StatusBadRequest, err
	}

	now := time.Now()
	if body.Expiry == (time.Time{}) {
		body.Expiry = now.Add(maxValidity / 2)
	}

	if body.Expiry.Before(now) || body.Expiry.Sub(now) > maxValidity {
		return http.StatusBadRequest, fmt.Errorf("token must be valid for max %s", maxValidity.String())
	}

	// TODO(flowlo): This will reject all scopes for now, as we are not using them.
	// As soon as we introduce scopes, this check must be rewritten accordingly.
	if len(body.Scopes) > 0 {
		return http.StatusBadRequest, fmt.Errorf("unknown scopes: %s", strings.Join(body.Scopes, ", "))
	}

	if len(body.Description) > 512 || len(body.Description) < 4 {
		return http.StatusBadRequest, errors.New("description has bad len")
	}

	// TODO(flowlo): Generate token with its own util.
	tokenValue, err := password.Generate(0)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	accessToken := model.AccessToken{
		Value:        string(tokenValue),
		Creation:     now,
		Modification: now,
		Expiry:       body.Expiry,
		Description:  body.Description,
		Scopes:       body.Scopes,
	}

	key, err := accessToken.SaveWithParent(ctx, p.UserKey)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	p.AccessToken = accessToken
	p.AccessTokenKey = key

	buf := new(bytes.Buffer)
	if err = gob.NewEncoder(buf).Encode(p); err != nil {
		return http.StatusInternalServerError, err
	}

	item := &memcache.Item{
		Key:   accessToken.Value,
		Value: buf.Bytes(),
	}

	if err = memcache.Set(ctx, item); err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(accessToken.Key(key))
	return
}
