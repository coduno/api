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

func Templates(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	taskKey, err := datastore.DecodeKey(mux.Vars(r)["key"])
	if err != nil {
		return http.StatusBadRequest, err
	}

	ls := r.URL.Query()["language"]

	var t model.Task
	if err = datastore.Get(ctx, taskKey, &t); err != nil {
		return http.StatusInternalServerError, nil
	}

	// TODO(flowlo): Use correct duration.
	expiry := time.Now().Add(time.Hour * 2)
	var urls []string

	expose := func(objs []model.StoredObject) error {
		for _, obj := range objs {
			u, err := util.Expose(obj.Bucket, obj.Name, expiry)
			if err != nil {
				return err
			}
			urls = append(urls, u)
		}
		return nil
	}

	if len(ls) == 0 {
		for _, objs := range t.Templates {
			if err := expose(objs); err != nil {
				return http.StatusInternalServerError, err
			}
		}
	} else {
		for _, l := range ls {
			if err := expose(t.Templates[l]); err != nil {
				return http.StatusInternalServerError, err
			}
		}
	}

	if len(urls) == 0 {
		w.Write([]byte("[]"))
	} else {
		json.NewEncoder(w).Encode(urls)
	}

	return http.StatusOK, nil
}
