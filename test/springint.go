package test

import (
	"io"

	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	RegisterTester(SpringInt, springInt)
}

func springInt(_ context.Context, t model.KeyedTest, sub model.KeyedSubmission, ball io.Reader) error {
	// TODO: use real context here if possible
	ctx := appengine.BackgroundContext()
	tr, err := runner.SpringInt(ctx, sub, ball)
	if err != nil {
		log.Debugf(ctx, "Spring Runner: getting runner results %+v", err)
		return err
	}

	if _, err := tr.Put(ctx, nil); err != nil {
		log.Debugf(ctx, "Spring Runner: error putting results in the datastore %+v", err)
		return err
	}

	log.Debugf(ctx, "Spring Runner: Marshalling")
	return marshalJSON(&sub, tr)
}
