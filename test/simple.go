package test

import (
	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

func init() {
	RegisterTester(Simple, simple)
}

func simple(ctx context.Context, params map[string]string, sub model.KeyedSubmission) error {
	log.Debugf(ctx, "Executing simple tester")
	stdout, stderr, err := runner.Simple(ctx, sub)
	log.Warningf(ctx, "%s %s %s", stdout, stderr, err)
	return err
}
