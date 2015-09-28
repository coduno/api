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
	RegisterTester(IO, iot)
}

func iot(ctx context.Context, t model.KeyedTest, sub model.KeyedSubmission, ball io.Reader) error {
	if err := checkIoParams(t.Params); err != nil {
		return err
	}

	ts, err := runner.IODiffRun(ctx, t, sub, ball)
	if err != nil {
		return err
	}

	// TODO(flowlo): Use a json.Encoder
	var body []byte
	if body, err = json.Marshal(ts); err != nil {
		return err
	}

	return ws.Write(sub.Key.Parent(), body)
}

func checkIoParams(params map[string]string) error {
	// TODO(victorbalan): check if len matches too
	for _, param := range [...]string{"bucket", "input", "output"} {
		if _, ok := params[param]; !ok {
			return ErrMissingParam(param)
		}
	}
	return nil
}
