package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"github.com/coduno/api/util/passenger"
	"github.com/gorilla/mux"

	"golang.org/x/net/context"
)

var pw = ""

func init() {
	if appengine.IsDevAppServer() {
		return
	}

	b, err := ioutil.ReadFile("pw")
	if err != nil {
		panic(err)
	}
	pw = strings.Trim(string(b), "\r\n ")
}

// PostSubmission creates a new submission.
func PostSubmission(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, errors.New("Unauthorized request")
	}
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, nil
	}

	submissionKind := r.URL.Query()["kind"][0]
	switch submissionKind {
	case "code":
		var body = struct {
			Task *datastore.Key
			Code,
			Language string
		}{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return http.StatusInternalServerError, err
		}
		var codeTask model.CodeTask
		if err = datastore.Get(ctx, body.Task, &codeTask); err != nil {
			return http.StatusInternalServerError, err
		}
		var response *http.Response
		if response, err = runOnDocker(w, codeTask, body.Language, body.Code); err != nil {
			return http.StatusInternalServerError, err
		}
		defer response.Body.Close()
		io.Copy(w, response.Body)
		// TODO(victorbalan): Process the engine response and create a submission.
		fallthrough
	case "question":
		return http.StatusOK, nil
	default:
		return http.StatusBadRequest, errors.New("Unknown submission kind.")
	}

	// TODO(victorbalan, flowlo): Create real submission from docker response.
	var submission model.Submission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		return http.StatusInternalServerError, err
	}

	resultKey, err := datastore.DecodeKey(mux.Vars(r)["key"])

	if !util.HasParent(p.UserKey, resultKey) {
		return http.StatusBadRequest, errors.New("Cannot submit answer for other users")
	}

	key, err := submission.SaveWithParent(ctx, resultKey)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	json.NewEncoder(w).Encode(submission.Key(key))
	return http.StatusCreated, nil
}

func runOnDocker(w http.ResponseWriter, task model.CodeTask, language, code string) (r *http.Response, err error) {
	var data = struct {
		Flags, Code, Language string
	}{
		task.Flags, code, language,
	}

	location := "https://engine.cod.uno/"
	if appengine.IsDevAppServer() {
		location = "http://localhost:8081/"
	}

	buf := new(bytes.Buffer)
	if err = json.NewEncoder(buf).Encode(data); err != nil {
		return
	}

	return http.Post(location+task.Runner, "application/json", buf)
}
