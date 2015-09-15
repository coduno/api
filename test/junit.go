package test

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"github.com/coduno/api/ws"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

func init() {
	RegisterTester(Junit, junit)
}

func junit(ctx context.Context, params map[string]string, sub model.KeyedSubmission) (err error) {
	log.Debugf(ctx, "Executing junit tester")
	if !checkJunitParams(params) {
		return errors.New("params missing")
	}
	for _, val := range strings.Split(params["tests"], ";") {
		var tr model.JunitTestResult
		if tr, err = runner.JUnit(ctx, val, sub); err != nil {
			return
		}
		if _, err = tr.Put(ctx, nil); err != nil {
			return
		}
		var body []byte
		if body, err = json.Marshal(tr); err != nil {
			return
		}
		if err = ws.Write(sub.Key, body); err != nil {
			return
		}
	}
	return
}

func checkJunitParams(params map[string]string) (ok bool) {
	if _, ok = params["tests"]; !ok {
		return
	}
	return true
}
