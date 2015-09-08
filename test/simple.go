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

func simple(ctx context.Context, params map[string]string, sub model.KeyedSubmission) error {
	log.Debugf(ctx, "Executing simple tester")
	str, err := runner.Simple(ctx, sub)
	log.Warningf(ctx, "%#v %s", str, err)

	// FIXME(victorbalan): Error handling
	j, _ := json.Marshal(str)
	return ws.Write(sub.Key, j)
}
