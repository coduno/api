package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/coduno/api/model"
	"golang.org/x/net/context"
)

func CreateCodeTask(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, nil
	}

	var body = struct {
		model.Task
		Flags     string
		Languages []string
		Runner    string
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return http.StatusBadRequest, err
	}

	task := model.CodeTask{
		Task:      body.Task,
		Flags:     body.Flags,
		Languages: body.Languages,
		Runner:    body.Runner,
	}

	key, err := task.Save(ctx, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(task.Key(key))
	return http.StatusOK, nil
}
