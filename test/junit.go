package test

import (
	"encoding/json"

	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"github.com/coduno/api/ws"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

func init() {
	RegisterTester(Junit, junit)
}

func junit(ctx context.Context, params map[string]string, sub model.KeyedSubmission) error {
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
	return ws.Write(sub.Key, j)
}
