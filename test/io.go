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
	RegisterTester(IO, io)
}

func io(ctx context.Context, t model.KeyedTest, sub model.KeyedSubmission) (err error) {
	log.Debugf(ctx, "Executing io tester")
	if !checkIoParams(t.Params) {
		return errors.New("params missing")
	}
	var ts model.TestStats
	if ts, err = runner.IODiffRun(ctx, t, sub); err != nil {
		return
	}

	var body []byte
	if body, err = json.Marshal(ts); err != nil {
		return err
	}
	return ws.Write(sub.Key, body)
}

func checkIoParams(params map[string]string) (ok bool) {
	if _, ok = params["bucket"]; !ok {
		return
	}
	// TODO(victorbalan): check if len matches too
	if _, ok = params["input"]; !ok {
		return
	}
	if _, ok = params["output"]; !ok {
		return
	}
	return true
}
