package test

import (
	"archive/tar"
	"io"
	"io/ioutil"
	"path"

	"google.golang.org/appengine"

	"github.com/coduno/api/model"
	"github.com/coduno/api/runner"
	"github.com/coduno/api/util"
	"golang.org/x/net/context"
)

func init() {
	RegisterTester(Junit, junit)
}

func junit(_ context.Context, t model.KeyedTest, sub model.KeyedSubmission, ball io.Reader) error {
	// TODO: use real context here if possible
	ctx := appengine.BackgroundContext()
	if _, ok := t.Params["test"]; !ok {
		return ErrMissingParam("test")
	}

	tests := model.StoredObject{
		Bucket: util.TestsBucket,
		Name:   t.Params["test"],
	}

	testStream := stream(ctx, tests)

	tr, err := runner.JUnit(ctx, testStream, ball)
	if err != nil {
		return err
	}

	if _, err := tr.PutWithParent(ctx, sub.Key); err != nil {
		return err
	}

	return marshalJSON(&sub, tr)
}

func stream(ctx context.Context, file model.StoredObject) io.ReadCloser {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		rc, err := util.Load(ctx, file.Bucket, file.Name)
		if err != nil {
			return
		}
		defer rc.Close()
		buf, err := ioutil.ReadAll(rc)
		if err != nil {
			return
		}
		w := tar.NewWriter(pw)
		defer w.Close()
		w.WriteHeader(&tar.Header{
			Name: path.Base(file.Name),
			Mode: 0600,
			Size: int64(len(buf)),
		})
		w.Write(buf)
	}()
	return pr
}
