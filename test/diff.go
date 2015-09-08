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

func diff(ctx context.Context, params map[string]string, sub model.KeyedSubmission) error {
	log.Debugf(ctx, "Executing diff tester")
	if !checkDiffParams(params) {
		return errors.New("params missing")
	}
	tr, err := runner.OutMatchDiffRun(ctx, params, sub)
	log.Warningf(ctx, "%#v %s", tr, err)

	// FIXME(victorbalan): Error handling
	j, _ := json.Marshal(tr)
	return ws.Write(sub.Key, j)
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
