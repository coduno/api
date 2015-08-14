package controllers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
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
	compute, err = url.Parse("https://" + credentials + "@compute.cod.uno")
	if err != nil {
		panic(err)
	}
}

// PostSubmission creates a new submission.
func PostSubmission(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, nil
	}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	resultKey, err := datastore.DecodeKey(mux.Vars(r)["resultKey"])

	if !util.HasParent(p.User, resultKey) {
		return http.StatusBadRequest, errors.New("cannot submit answer for other users")
	}

	taskKey, err := datastore.DecodeKey(mux.Vars(r)["taskKey"])
	// Note: When more task kinds are added, see controllers.CreateFinalResult.
	switch taskKey.Kind() {
	case model.CodeTaskKind:
		return runner.HandleCodeSubmission(ctx, w, r, resultKey, taskKey)
	// TODO(victorbalan, flowlo): Use correct kind when possible.
	case "QuestionTask":
		return http.StatusInternalServerError, errors.New("question submissions are not yet implemented")
	default:
		return http.StatusBadRequest, errors.New("Unknown submission kind.")
	}
}

// FinalSubmission makes the last submission final.
func FinalSubmission(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, nil
	}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	var resultKey *datastore.Key
	if resultKey, err = datastore.DecodeKey(mux.Vars(r)["resultKey"]); err != nil {
		return http.StatusInternalServerError, err
	}

	if !util.HasParent(p.User, resultKey) {
		return http.StatusBadRequest, errors.New("cannot submit answer for other users")
	}

	var index int
	if index, err = strconv.Atoi(mux.Vars(r)["index"]); err != nil {
		return http.StatusInternalServerError, err
	}

	var submissionKey *datastore.Key
	if submissionKey, err = datastore.DecodeKey(mux.Vars(r)["submissionKey"]); err != nil {
		return http.StatusInternalServerError, err
	}

	var result model.Result
	if err = datastore.Get(ctx, resultKey, &result); err != nil {
		return http.StatusInternalServerError, err
	}

	result.FinalSubmissions[index] = submissionKey

	if _, err = result.Save(ctx, resultKey); err != nil {
		return http.StatusInternalServerError, err
	}
	w.Write([]byte("OK"))
	return
}
