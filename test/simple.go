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
	RegisterTester(Simple, simple)
}

func simple(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string, sub model.KeyedSubmission) error {
	log.Debugf(ctx, "Executing simple tester")
	stdout, stderr, err := runner.Simple(ctx, sub)
	log.Warningf(ctx, "%s %s %s", stdout, stderr, err)

	j, _ := json.Marshal(struct {
		Stdout string
		Stderr string
	}{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	})
	_, err = w.Write(j)
	return err
}
