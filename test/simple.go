package test

import (
	"io"

	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"golang.org/x/net/context"
)

func init() {
	RegisterTester(Simple, simple)
}

func simple(ctx context.Context, t model.Test, sub model.Submission, ball io.Reader) (err error) {
	str, err := runner.Simple(ctx, sub, ball)
	if err != nil {
		return
	}

	return marshalJSON(&sub, str)
}
