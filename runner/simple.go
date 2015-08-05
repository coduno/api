package runner

import (
	"encoding/json"
	"net/http"

	"github.com/coduno/api/model"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// SimpleRunner is the runner used for a simple run.
type SimpleRunner struct {
	Submission model.CodeSubmission
}

func (sr *SimpleRunner) Run(ctx context.Context, w http.ResponseWriter, r *http.Request, codeTask model.CodeTask, resultKey *datastore.Key) (status int, err error) {
	response, err := sr.Start(ctx, w, r, codeTask)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	sr.Handle(ctx, w, response, resultKey)
	return
}

// Start function for a simple run.
func (sr *SimpleRunner) Start(ctx context.Context, w http.ResponseWriter, r *http.Request, codeTask model.CodeTask) (response *http.Response, err error) {
	if err = json.NewDecoder(r.Body).Decode(&sr.Submission); err != nil {
		return
	}
	return run(codeTask, sr.Submission.Language, sr.Submission.Code)
}

// Handle function for a simple run.
func (sr *SimpleRunner) Handle(ctx context.Context, w http.ResponseWriter, response *http.Response, resultKey *datastore.Key) {
	json.NewDecoder(response.Body).Decode(&sr.Submission)

	key, err := sr.Submission.SaveWithParent(ctx, resultKey)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(sr.Submission.Key(key))
}
