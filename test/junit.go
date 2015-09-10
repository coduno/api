package test

import (
	"encoding/json"
	"net/http"

	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

func init() {
	RegisterTester(Junit, junit)
}

func junit(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string, sub model.KeyedSubmission) error {
	log.Debugf(ctx, "Executing junit tester")
	stdout, stderr, utr, err := runner.JUnit(ctx, params, sub)
	log.Warningf(ctx, "%s %s %s", stdout, stderr, err)

	j, _ := json.Marshal(struct {
		Stdout  string
		Stderr  string
		Results []model.UnitTestResults
	}{
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
		Results: utr,
	})
	_, err = w.Write(j)
	return err
}
