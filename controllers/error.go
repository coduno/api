package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"github.com/coduno/api/util/passenger"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func Error(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if util.CheckMethod(r, "GET") {
		return http.StatusMethodNotAllowed, nil
	}

	var body = struct {
		Status                  int
		User                    *datastore.Key
		Key, Description, Route string
	}{}

	if p, ok := passenger.FromContext(ctx); ok {
		body.User = p.UserKey
	}

	err = json.NewDecoder(r.Body).Decode(&body)

	userError := model.Error{
		Status:      body.Status,
		User:        body.User,
		Description: body.Description,
		Route:       body.Route}

	if body.Key != "" {
		key, err := datastore.DecodeKey(body.Key)
		if err != nil {
			return http.StatusInternalServerError, nil
		}
		userError.Save(ctx, key)
	} else {
		userError.Save(ctx)
	}
	return http.StatusOK, nil
}
