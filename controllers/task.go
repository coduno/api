package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/coduno/api/db"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

func init() {
	router.Handle("/tasks/{id}/tests", ContextHandlerFunc(TestsByTaskKey))
}

func TestsByTaskKey(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	taskId, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tests, err := db.LoadTestsForTask(taskId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(tests)
	return http.StatusOK, nil
}
