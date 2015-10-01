package test

import (
	"encoding/json"
	"io"
	"strconv"

	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"github.com/coduno/api/util"
	"github.com/coduno/api/ws"
	"golang.org/x/net/context"
)

func init() {
	RegisterTester(CoderJunit, coderJunit)
}

func coderJunit(ctx context.Context, t model.KeyedTest, sub model.KeyedSubmission, ball io.Reader) (err error) {
	if _, ok := t.Params["code"]; !ok {
		return ErrMissingParam("code")
	}

	code := model.StoredObject{
		Bucket: util.TestsBucket,
		Name:   t.Params["code"],
	}
	codeStream := stream(ctx, code)
	var tr *model.JunitTestResult
	if tr, err = runner.JUnit(ctx, ball, codeStream); err != nil {
		return
	}

	shouldFail, err := strconv.ParseBool(t.Params["shouldFail"])
	if err != nil {
		return
	}

	ctr := model.CoderJunitTestResult{
		JunitTestResult: *tr,
		ShouldFail:      shouldFail,
	}
	if _, err = ctr.Put(ctx, nil); err != nil {
		return
	}

	var body []byte
	if body, err = json.Marshal(ctr); err != nil {
		return
	}
	return ws.Write(sub.Key.Parent(), body)
}
