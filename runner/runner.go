package runner

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/coduno/api/model"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
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

// Runner is the general runner interface wich will start a docker run and
// save the results.
type Runner interface {
	Run(ctx context.Context, w http.ResponseWriter, r *http.Request, codeTask model.CodeTask, resultKey *datastore.Key) (status int, err error)
}

// HandleCodeSubmission starts the correct runner depending on the codeTask runner.
func HandleCodeSubmission(ctx context.Context, w http.ResponseWriter, r *http.Request, resultKey, taskKey *datastore.Key) (status int, err error) {
	var codeTask model.CodeTask
	if err = datastore.Get(ctx, taskKey, &codeTask); err != nil {
		return http.StatusInternalServerError, err
	}
	var runner Runner
	switch codeTask.Runner {
	case "simple":
		submission := model.CodeSubmission{Submission: model.Submission{Task: taskKey}}
		runner = &SimpleRunner{Submission: submission}
	case "javaut":
		submission := model.JunitSubmission{Submission: model.Submission{Task: taskKey}}
		runner = &JunitRunner{Submission: submission}
	default:
		return http.StatusBadRequest, nil
	}
	return runner.Run(ctx, w, r, codeTask, resultKey)
}

func run(codeTask model.CodeTask, language, code string) (r *http.Response, err error) {
	var data = struct {
		Flags, Code, Language string
	}{
		codeTask.Flags, code, language,
	}

	buf := new(bytes.Buffer)
	if err = json.NewEncoder(buf).Encode(data); err != nil {
		return
	}

	return http.Post(compute.String()+"/"+codeTask.Runner, "application/json;charset=utf-8", buf)
}
