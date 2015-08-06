package runner

import (
	"encoding/json"
	"net/http"

	"github.com/coduno/api/model"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// DiffRunner is the runner used for a diff run.
type DiffRunner struct {
	Submission model.DiffSubmission
}

func (dr *DiffRunner) Run(ctx context.Context, w http.ResponseWriter, r *http.Request, codeTask model.CodeTask, resultKey *datastore.Key) (status int, err error) {
	response, err := dr.Start(ctx, w, r, codeTask)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	dr.Handle(ctx, w, response, resultKey)
	return
}

// Start function for a diff run.
func (dr *DiffRunner) Start(ctx context.Context, w http.ResponseWriter, r *http.Request, codeTask model.CodeTask) (response *http.Response, err error) {
	if err = json.NewDecoder(r.Body).Decode(&dr.Submission); err != nil {
		return
	}
	return run(codeTask, dr.Submission.Language, dr.Submission.Code)
}

// Handle function for a diff run.
func (dr *DiffRunner) Handle(ctx context.Context, w http.ResponseWriter, response *http.Response, resultKey *datastore.Key) {
	json.NewDecoder(response.Body).Decode(&dr.Submission)

	key, err := dr.Submission.SaveWithParent(ctx, resultKey)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(dr.Submission.Key(key))
}
