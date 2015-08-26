package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"github.com/coduno/api/util/passenger"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// TaskByKey loads a task by key.
func TaskByKey(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	taskKey, err := datastore.DecodeKey(mux.Vars(r)["key"])
	if err != nil {
		return http.StatusInternalServerError, err
	}

	var u model.User
	if err = datastore.Get(ctx, p.User, &u); err != nil {
		return http.StatusInternalServerError, nil
	}

	// User is a coder
	if u.Company == nil {
		rk, err := datastore.DecodeKey(r.URL.Query()["result"][0])
		if err != nil {
			return http.StatusInternalServerError, err
		}

		if !util.HasParent(p.User, rk) {
			return http.StatusUnauthorized, nil
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

	var task model.Task
	if err = datastore.Get(ctx, taskKey, &task); err != nil {
		return http.StatusInternalServerError, err
	}
	json.NewEncoder(w).Encode(task.Key(taskKey))
	return http.StatusOK, nil
}

func Tasks(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	var u model.User
	if err = datastore.Get(ctx, p.User, &u); err != nil {
		return http.StatusInternalServerError, nil
	}

	// User is a coder
	if u.Company == nil {
		return http.StatusUnauthorized, nil
	}

	var tasks model.Tasks
	taskKeys, err := model.NewQueryForTask().
		GetAll(ctx, &tasks)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(tasks.Key(taskKeys))
	return http.StatusOK, nil
}
