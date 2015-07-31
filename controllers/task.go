package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/coduno/app/model"
	"github.com/coduno/app/util"
	"github.com/coduno/app/util/passenger"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// TaskByKey loads a task by key.
func TaskByKey(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, errors.New("Unauthorized request")
	}

	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	taskKey, err := datastore.DecodeKey(mux.Vars(r)["key"])
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// User is a coder
	if p.UserKey.Parent() == nil {
		rk, err := datastore.DecodeKey(r.URL.Query()["result"][0])
		if err != nil {
			return http.StatusInternalServerError, err
		}

		if !util.HasParent(p.UserKey, rk) {
			return http.StatusUnauthorized, errors.New("Unauthorized")
		}

		var result model.Result
		if err = datastore.Get(ctx, rk, &result); err != nil {
			return http.StatusInternalServerError, err
		}

		var challenge model.Challenge
		if err = datastore.Get(ctx, result.Challenge, &challenge); err != nil {
			return http.StatusInternalServerError, err
		}

		emptyTime := time.Time{}
		updateResult := false
		for i, val := range challenge.Tasks {
			if taskKey.Equal(val) && result.StartTimes[i] == emptyTime {
				result.StartTimes[i] = time.Now()
				updateResult = true
				break
			}
		}
		if updateResult {
			if _, err = result.Save(ctx, rk); err != nil {
				return http.StatusInternalServerError, err
			}
		}
	}

	switch taskKey.Kind() {
	case "codeTasks":
		var task model.CodeTask
		if err = datastore.Get(ctx, taskKey, &task); err != nil {
			return http.StatusInternalServerError, err
		}
		json.NewEncoder(w).Encode(task.Key(taskKey))
		return http.StatusOK, nil
	default:
		return http.StatusInternalServerError, errors.New("Not yet impl")
	}
}
