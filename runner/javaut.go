package runner

import (
	"github.com/coduno/api/model"
	"golang.org/x/net/context"
)

func junitRunner(ctx context.Context, test *model.Test, sub model.KeyedSubmission) error {
	// response, err := jr.Start(ctx, w, r, codeTask)
	// if err != nil {
	// 	return http.StatusInternalServerError, err
	// }
	//jr.Handle(ctx, w, response, resultKey)
	return nil
}

//
// // Start function for a java unit test run.
// func (jr *JunitRunner) Start(ctx context.Context, w http.ResponseWriter, r *http.Request, codeTask model.CodeTask) (response *http.Response, err error) {
// 	if err = json.NewDecoder(r.Body).Decode(&jr.Submission); err != nil {
// 		return
// 	}
// 	// TODO(victorbalan): pass the correct language when we have different JUnit tests
// 	return run(codeTask, "javaut", jr.Submission.Code)
// }
//
// // Handle function for a JUnit test run.
// func (jr *JunitRunner) Handle(ctx context.Context, w http.ResponseWriter, response *http.Response, resultKey *datastore.Key) {
// 	json.NewDecoder(response.Body).Decode(&jr.Submission)
//
// 	key, err := jr.Submission.SaveWithParent(ctx, resultKey)
// 	if err != nil {
// 		return
// 	}
// 	json.NewEncoder(w).Encode(jr.Submission.Key(key))
// }
