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

func diff(ctx context.Context, params map[string]string, sub model.KeyedSubmission) (err error) {
	log.Debugf(ctx, "Executing diff tester")
	if !checkDiffParams(params) {
		return errors.New("params missing")
	}
	var tr model.DiffTestResult
	if tr, err = runner.OutMatchDiffRun(ctx, params, sub); err != nil {
		return
	}
	if _, err = tr.Put(ctx, nil); err != nil {
		return
	}
	var body []byte
	if body, err = json.Marshal(tr); err != nil {
		return
	}
	return ws.Write(sub.Key, body)
}

func checkDiffParams(params map[string]string) (ok bool) {
	if _, ok = params["bucket"]; !ok {
		return
	}
	if _, ok = params["tests"]; !ok {
		return
	}
	return true
}
