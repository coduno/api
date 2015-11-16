package test

import (
	"io"

	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"golang.org/x/net/context"
)

func init() {
	RegisterTester(IO, iot)
}

func iot(ctx context.Context, t model.Test, sub model.Submission, ball io.Reader) error {
	if err := checkIoParams(t.Params); err != nil {
		return err
	}

	ts, err := runner.IODiffRun(ctx, t, sub, ball)
	if err != nil {
		return err
	}

	return marshalJSON(&sub, ts)
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
