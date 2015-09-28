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
	RegisterTester(Simple, simple)
}

func simple(ctx context.Context, t model.KeyedTest, sub model.KeyedSubmission, ball io.Reader) (err error) {
	str, err := runner.Simple(ctx, sub, ball)
	if err != nil {
		return
	}

	// TODO(flowlo): Use a json.Encoder
	var body []byte
	if body, err = json.Marshal(str); err != nil {
		return
	}
	return ws.Write(sub.Key.Parent(), body)
}
