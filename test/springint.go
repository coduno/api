package test

import (
	"io"

	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"golang.org/x/net/context"
)

func init() {
	RegisterTester(SpringInt, springInt)
}

func springInt(ctx context.Context, t model.KeyedTest, sub model.KeyedSubmission, ball io.Reader) error {
	tr, err := runner.SpringInt(ctx, sub, ball)
	if err != nil {
		return err
	}

	if _, err := tr.Put(ctx, nil); err != nil {
		return err
	}

	return marshalJSON(&sub, tr)
}
