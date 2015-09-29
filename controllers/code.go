package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"github.com/gorilla/mux"

	"golang.org/x/net/context"
)

func init() {
	router.Handle("/tasks/{key}/templates", ContextHandlerFunc(Templates))
}

// Template serves the contents of a static file to a client.
// TODO(flowlo, victorbalan): Decide where the templates will be stored.
func Templates(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	taskKey, err := datastore.DecodeKey(mux.Vars(r)["key"])
	if err != nil {
		return http.StatusBadRequest, err
	}

	var t model.Task
	if err = datastore.Get(ctx, taskKey, &t); err != nil {
		return http.StatusInternalServerError, nil
	}

	// TODO(flowlo): Use correct duration.
	expiry := time.Now().Add(time.Hour * 2)

	urls := make([]string, 0, len(t.Templates))
	for _, template := range t.Templates {
		u, err := util.Expose(template.Bucket, template.Name, expiry)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		urls = append(urls, u)
	}

	json.NewEncoder(w).Encode(urls)
	return http.StatusOK, nil
}
