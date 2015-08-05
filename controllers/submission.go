package controllers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/coduno/api/runner"
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
	compute, err = url.Parse("https://" + credentials + "@git.cod.uno")
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

	resultKey, err := datastore.DecodeKey(mux.Vars(r)["resultKey"])

	if !util.HasParent(p.UserKey, resultKey) {
		return http.StatusBadRequest, errors.New("Cannot submit answer for other users")
	}

	taskKey, err := datastore.DecodeKey(mux.Vars(r)["taskKey"])

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
