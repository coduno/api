package controllers

import (
	"net/http"
	"time"

	"github.com/coduno/engine/model"
	"github.com/coduno/engine/passenger"
	"github.com/coduno/engine/util"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// GetTaskByKey loads a task by key
func GetTaskByKey(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	p, ok := passenger.FromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if !util.CheckMethod(w, r, "GET") {
		return
	}
	taskKey, err := datastore.DecodeKey(mux.Vars(r)["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// User is a coder
	if p.UserKey.Parent() == nil {
		rk, err := datastore.DecodeKey(r.URL.Query()["result"][0])
		if err != nil {
			http.Error(w, "Key decoding err "+err.Error(), http.StatusInternalServerError)
			return
		}
		if !util.HasParent(p.UserKey, rk) {
			http.Error(w, "Unauthorized to see the task", http.StatusUnauthorized)
			return
		}
		var result model.Result
		err = datastore.Get(ctx, rk, &result)
		if err != nil {
			http.Error(w, "Datastore err "+err.Error(), http.StatusInternalServerError)
			return
		}
		var challenge model.Challenge
		err = datastore.Get(ctx, result.Challenge, &challenge)
		if err != nil {
			http.Error(w, "Datastore err "+err.Error(), http.StatusInternalServerError)
			return
		}

		emptyTime := time.Unix(0, 0)
		updateResult := false
		for i, val := range challenge.Tasks {
			if taskKey.Equal(val) && result.StartTimes[i] == emptyTime {
				result.StartTimes[i] = time.Now()
				updateResult = true
				break
			}
		}
		if updateResult {
			_, err = result.Save(ctx, rk)
			if err != nil {
				http.Error(w, "Datastore err "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

	}

	var task model.Task
	err = datastore.Get(ctx, taskKey, &task)
	if err != nil {
		http.Error(w, "Datastore err"+err.Error(), http.StatusInternalServerError)
		return
	}
	task.Write(w, taskKey)
}
