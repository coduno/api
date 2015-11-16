package test

import (
	"io"

	"github.com/coduno/api/db"
	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"golang.org/x/net/context"
)

func init() {
	RegisterTester(Junit, junit)
}

func junit(ctx context.Context, t model.Test, sub model.Submission, ball io.Reader) error {
	if _, ok := t.Params["test"]; !ok {
		return ErrMissingParam("test")
	}

	testStream := db.LoadFile(t.Params["test"])

	tr, err := runner.JUnit(ctx, testStream, ball)
	if err != nil {
		return err
	}

	// if _, err := tr.PutWithParent(ctx, sub.Key); err != nil {
	// 	return err
	// }

	return marshalJSON(&sub, tr)
}
