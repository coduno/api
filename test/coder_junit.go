package test

import (
	"io"
	"strconv"

	"google.golang.org/appengine/log"

	"github.com/coduno/api/db"
	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"golang.org/x/net/context"
)

func init() {
	RegisterTester(CoderJunit, coderJunit)
}

func coderJunit(ctx context.Context, t model.Test, sub model.Submission, ball io.Reader) (err error) {
	if _, ok := t.Params["code"]; !ok {
		log.Debugf(ctx, "JUnit Coder Runner: missing param")
		return ErrMissingParam("code")
	}

	codeStream := db.LoadFile(t.Params["code"])
	var tr *model.JunitTestResult
	if tr, err = runner.JUnit(ctx, ball, codeStream); err != nil {
		log.Debugf(ctx, "JUnit Coder Runner: should fail %+v", err)
		return
	}

	shouldFail, err := strconv.ParseBool(t.Params["shouldFail"])
	if err != nil {
		log.Debugf(ctx, "JUnit Coder Runner: should fail %+v", err)
		return
	}

	ctr := model.CoderJunitTestResult{
		JunitTestResult: *tr,
		ShouldFail:      shouldFail,
	}

	// if _, err = ctr.Put(ctx, nil); err != nil {
	// 	log.Debugf(ctx, "JUnit Coder Runner: store %+v", err)
	// 	return
	// }

	return marshalJSON(&sub, processResult(t, ctr))
}

func processResult(t model.Test, result model.CoderJunitTestResult) (ts model.TestStats) {
	ts = model.TestStats{
		Stdout: result.Stdout,
		Test:   t.ID,
	}

	if result.Results.Tests == 0 {
		ts.Stderr = result.Stderr
		ts.Failed = true
		return
	}

	ts.Failed = result.ShouldFail == (result.Results.Failures == 0)
	return
}
