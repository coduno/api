package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"github.com/coduno/api/util/passenger"
	"github.com/gorilla/mux"

	"golang.org/x/net/context"
)

var compute *url.URL

func init() {
	var err error
	if appengine.IsDevAppServer() {
		compute, err = url.Parse("http://localhost:8081")
		if err != nil {
			panic(err)
		}
		return
	}

	b, err := ioutil.ReadFile("credentials")
	if err != nil {
		panic(err)
	}

	credentials := strings.Trim(string(b), "\r\n ")
	compute, err = url.Parse("https://" + credentials + "git.cod.uno")
	if err != nil {
		panic(err)
	}
}

// PostSubmission creates a new submission.
func PostSubmission(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, nil
	}

	resultKey, err := datastore.DecodeKey(mux.Vars(r)["key"])

	if !util.HasParent(p.UserKey, resultKey) {
		return http.StatusBadRequest, errors.New("Cannot submit answer for other users")
	}

	var capacity int
	capacity, err = strconv.Atoi(r.Header["Content-Length"][0])
	body := make([]byte, 0, capacity)
	body, err = ioutil.ReadAll(r.Body)

	var submission model.Submission
	json.Unmarshal(body, &submission)

	switch submission.Task.Kind() {
	case "codeTasks":
		return codeSubmission(ctx, w, body, resultKey)
	case "questionTasks":
		return http.StatusInternalServerError, errors.New("question submissions are not yet implemented")
	default:
		return http.StatusBadRequest, errors.New("Unknown submission kind.")
	}
}

func codeSubmission(ctx context.Context, w http.ResponseWriter, body []byte, resultKey *datastore.Key) (status int, err error) {
	var submission model.CodeSubmission
	if err := json.Unmarshal(body, &submission); err != nil {
		return http.StatusInternalServerError, err
	}

	var codeTask model.CodeTask
	if err = datastore.Get(ctx, submission.Task, &codeTask); err != nil {
		return http.StatusInternalServerError, err
	}

	var response *http.Response
	if response, err = runOnDocker(w, codeTask, submission.Language, submission.Code); err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewDecoder(response.Body).Decode(&submission)

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

	buf := new(bytes.Buffer)
	if err = json.NewEncoder(buf).Encode(data); err != nil {
		return
	}

	return http.Post(compute.String()+"/"+task.Runner, "application/json;charset=utf-8", buf)
}
