package runner

import (
	"encoding/json"
	"net/http"

	"github.com/coduno/api/model"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// JunitRunner is the runner used for a JUnit test run.
type JunitRunner struct {
	Submission model.JunitSubmission
}

func (jr *JunitRunner) Run(ctx context.Context, w http.ResponseWriter, r *http.Request, codeTask model.CodeTask, resultKey *datastore.Key) (status int, err error) {
	response, err := jr.Start(ctx, w, r, codeTask)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	jr.Handle(ctx, w, response, resultKey)
	return
}

// Start function for a java unit test run.
func (jr *JunitRunner) Start(ctx context.Context, w http.ResponseWriter, r *http.Request, codeTask model.CodeTask) (response *http.Response, err error) {
	if err = json.NewDecoder(r.Body).Decode(&jr.Submission); err != nil {
		return
	}
	// TODO(victorbalan): pass the correct language when we have different JUnit tests
	return run(codeTask, "javaut", jr.Submission.Code)
}

// Handle function for a JUnit test run.
func (jr *JunitRunner) Handle(ctx context.Context, w http.ResponseWriter, response *http.Response, resultKey *datastore.Key) {
	json.NewDecoder(response.Body).Decode(&jr.Submission)

	key, err := jr.Submission.SaveWithParent(ctx, resultKey)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(jr.Submission.Key(key))
}
