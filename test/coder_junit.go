package test

import (
	"io"
	"strconv"

	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"github.com/coduno/api/util"
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

	return marshalJSON(&sub, processResult(t, ctr))
}

func processResult(t model.KeyedTest, result model.CoderJunitTestResult) (ts model.TestStats) {
	ts = model.TestStats{
		Stdout: result.Stdout,
		Test:   t.Key,
	}

	if result.Results.Tests == 0 {
		ts.Stderr = result.Stderr
		ts.Failed = true
		return
	}
	var failed bool
	if result.ShouldFail && result.Results.Failures == 0 {
		failed = true
	} else if result.ShouldFail && result.Results.Failures > 0 {
		failed = false
	} else if !result.ShouldFail && result.Results.Failures == 0 {
		failed = false
	} else {
		failed = true
	}
	ts.Failed = failed
	return
}
