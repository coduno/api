package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/coduno/app/model"
	"github.com/coduno/app/util"
	"github.com/coduno/engine/passenger"
	"github.com/gorilla/mux"

	"golang.org/x/net/context"
)

// CodeData is a hack
// TODO(victorbalan): Remove hack when we refactor the engine
type CodeData struct {
	Flags    string
	Code     string
	Runner   string
	Language string
}

// SubmissionData is a hack
// TODO(victorbalan): Remove hack when we refactor the engine
type SubmissionData struct {
	Task *datastore.Key
	Code,
	Language string
}

// PostSubmission creates a new submission
func PostSubmission(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, errors.New("Unauthorized request")
	}
	if err = util.CheckMethod(r, "POST"); err != nil {
		return http.StatusMethodNotAllowed, err
	}
	var submission model.Submission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		return http.StatusInternalServerError, err
	}

	submissionKind := r.URL.Query()["kind"][0]
	switch submissionKind {
	case "code":
		// TODO(victorbalan): Load body in separate struct and not in submission
		var submissionData SubmissionData
		if err := json.NewDecoder(r.Body).Decode(&submissionData); err != nil {
			return http.StatusInternalServerError, err
		}
		var codeTask model.CodeTask
		if err = datastore.Get(ctx, submissionData.Task, &codeTask); err != nil {
			return http.StatusInternalServerError, err
		}
		runOnDocker(w, codeTask, submissionData)
		// TODO(victorbalan): Process the engine response and create a submission.
		fallthrough
	case "question":
		return http.StatusOK, nil
	default:
		return http.StatusBadRequest, errors.New("Unknown submission kind.")
	}
	resultKey, err := datastore.DecodeKey(mux.Vars(r)["id"])

	if !util.HasParent(p.UserKey, resultKey) {
		return http.StatusBadRequest, errors.New("Cannot submit answer for other users")
	}

	key, err := submission.SaveWithParent(ctx, resultKey)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	submission.Write(w, key)
	return http.StatusCreated, nil
}

func runOnDocker(w http.ResponseWriter, task model.CodeTask, sd SubmissionData) {
	data := CodeData{
		Flags:    task.Flags,
		Runner:   task.Runner,
		Code:     sd.Code,
		Language: sd.Language,
	}
	engine := "https://engine.cod.uno"
	if appengine.IsDevAppServer() {
		engine = "http://localhost:8081"
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		return
	}
	res, _ := http.Post(engine, "json", buf)
	io.Copy(w, res.Body)
}
