package test

import (
	"io"

	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"golang.org/x/net/context"
)

func init() {
	RegisterTester(Diff, diff)
}

func diff(ctx context.Context, t model.Test, sub model.Submission, ball io.Reader) (err error) {
	if _, ok := t.Params["tests"]; !ok {
		return ErrMissingParam("tests")
	}

	ts, err := runner.OutMatchDiffRun(ctx, t, sub, ball)
	if err != nil {
		return
	}

	return marshalJSON(&sub, ts)
}
