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
	RegisterTester(Diff, diff)
}

func diff(ctx context.Context, t model.KeyedTest, sub model.KeyedSubmission, ball io.Reader) (err error) {
	if _, ok := t.Params["tests"]; !ok {
		return ErrMissingParam("tests")
	}

	ts, err := runner.OutMatchDiffRun(ctx, t, sub, ball)
	if err != nil {
		return
	}

	// TODO(flowlo): Use a json.Encoder
	var body []byte
	if body, err = json.Marshal(ts); err != nil {
		return err
	}
	return ws.Write(sub.Key.Parent(), body)
}
