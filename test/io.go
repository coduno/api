package test

import (
	"errors"
	"strings"

	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

func init() {
	RegisterTester(IO, io)
}

func io(ctx context.Context, params map[string]string, sub model.KeyedSubmission) error {
	log.Debugf(ctx, "Executing io tester")
	if !checkIoParams(params) {
		return errors.New("params missing")
	}
	in := strings.Split(params["input"], " ")
	out := strings.Split(params["input"], " ")
	for i := range in {
		tr, err := runner.IODiffRun(ctx, in[i], out[i], sub)
		log.Warningf(ctx, "%#v %s", tr, err)
		// TODO(victorbalan, flowlo): pass back the results through ws
	}
	return nil
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
