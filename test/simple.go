package test

import (
	"encoding/json"

	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"github.com/coduno/api/ws"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

func init() {
	RegisterTester(Simple, simple)
}

func simple(ctx context.Context, t model.KeyedTest, sub model.KeyedSubmission) (err error) {
	log.Debugf(ctx, "Executing simple tester")
	var str model.SimpleTestResult
	if str, err = runner.Simple(ctx, sub); err != nil {
		return
	}

	var body []byte
	if body, err = json.Marshal(str); err != nil {
		return
	}
	return ws.Write(sub.Key.Parent(), body)
}
