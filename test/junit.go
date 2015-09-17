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
	RegisterTester(Junit, junit)
}

func junit(ctx context.Context, t model.KeyedTest, sub model.KeyedSubmission) (err error) {
	log.Debugf(ctx, "Executing junit tester")
	if !checkJunitParams(t.Params) {
		return errors.New("params missing")
	}
	var tr model.JunitTestResult
	if tr, err = runner.JUnit(ctx, t.Params["test"], sub); err != nil {
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

func checkJunitParams(params map[string]string) (ok bool) {
	if _, ok = params["test"]; !ok {
		return
	}
	return true
}
