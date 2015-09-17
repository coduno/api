package test

import (
	"encoding/json"
	"errors"

	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"github.com/coduno/api/ws"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

func init() {
	RegisterTester(Diff, diff)
}

func diff(ctx context.Context, t model.KeyedTest, sub model.KeyedSubmission) (err error) {
	log.Debugf(ctx, "Executing diff tester")
	if !checkDiffParams(t.Params) {
		return errors.New("params missing")
	}
	var ts model.TestStats
	if ts, err = runner.OutMatchDiffRun(ctx, t, sub); err != nil {
		return
	}

	var body []byte
	if body, err = json.Marshal(ts); err != nil {
		return err
	}
	return ws.Write(sub.Key, body)
}

func checkDiffParams(params map[string]string) (ok bool) {
	if _, ok = params["tests"]; !ok {
		return
	}
	return true
}
