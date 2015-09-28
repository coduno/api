package test

import (
	"encoding/json"
	"io"

	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"github.com/coduno/api/ws"
	"golang.org/x/net/context"
)

func init() {
	RegisterTester(Junit, junit)
}

func junit(ctx context.Context, t model.KeyedTest, sub model.KeyedSubmission, ball io.Reader) error {
	if _, ok := t.Params["test"]; !ok {
		return ErrMissingParam("test")
	}

	tr, err := runner.JUnit(ctx, t.Params["test"], sub, ball)
	if err != nil {
		return err
	}

	if _, err := tr.Put(ctx, nil); err != nil {
		return err
	}

	// TODO(flowlo): Use a json.Encoder
	body, err := json.Marshal(tr)
	if err != nil {
		return err
	}

	return ws.Write(sub.Key.Parent(), body)
}
