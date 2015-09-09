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

	if len(r.URL.Query()["result"]) > 0 {
		rk, err := datastore.DecodeKey(r.URL.Query()["result"][0])
		if err != nil {
			return http.StatusInternalServerError, err
		}

		if util.HasParent(p.User, rk) {
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
				if taskKey.Equal(val) {
					if result.StartTimes[i].Equal(emptyTime) {
						result.StartTimes[i] = time.Now()
						updateResult = true
						break
					}
				}
			}
			if updateResult {
				if _, err = result.Put(ctx, rk); err != nil {
					return http.StatusInternalServerError, err
				}
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
	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	var u model.User
	if err = datastore.Get(ctx, p.User, &u); err != nil {
		return http.StatusInternalServerError, err
	}

	// User is a coder
	if u.Company == nil {
		return http.StatusUnauthorized, nil
	}

	switch r.Method {
	case "GET":
		return getAllTasks(ctx, w, r)
	case "POST":
		return createTask(ctx, w, r)
	default:
		return http.StatusMethodNotAllowed, nil
	}
}

func createTask(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	var task model.Task

	if err = json.NewDecoder(r.Body).Decode(&task); err != nil {
		return http.StatusBadRequest, err
	}

	var key *datastore.Key
	if key, err = task.Put(ctx, nil); err != nil {
		return http.StatusInternalServerError, err
	}
	json.NewEncoder(w).Encode(task.Key(key))
	return http.StatusOK, nil
}

func getAllTasks(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	var tasks model.Tasks
	taskKeys, err := model.NewQueryForTask().
		GetAll(ctx, &tasks)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(tasks.Key(taskKeys))
	return http.StatusOK, nil
}
