package util

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
	"google.golang.org/cloud/storage"
)

var FileNames = map[string]string{
	"py":    "app.py",
	"c":     "app.c",
	"cpp":   "app.cpp",
	"java":  "Application.java",
	"robot": "robot.json",
}

const (
	TemplateBucket = "coduno-templates"
	TestsBucket    = "coduno-tests"
	// TODO(victorbalan): Add param in the test struct to not hardcode
	// the result file name.
	JUnitResultsPath = "/run/build/test-results/TEST-Tests.xml"
)

var jwtc *jwt.Config

func init() {
	if raw, err := ioutil.ReadFile("service-account.json"); err == nil {
		jwtc, _ = google.JWTConfigFromJSON(raw, storage.ScopeFullControl)
	}
}

func SubmissionBucket() string {
	if appengine.IsDevAppServer() {
		return "coduno-dev"
	}
	return "coduno"
}

func Load(ctx context.Context, bucket, name string) (io.ReadCloser, error) {
	r, err := fromCache(ctx, bucket, name)
	if err == nil {
		return r, nil
	}
	if err != memcache.ErrCacheMiss {
		return nil, err
	}
	return fromStorage(ctx, bucket, name)
}

type cachingCloser struct {
	ctx context.Context
	key string
	rc  io.ReadCloser
	buf *bytes.Buffer
}

func (c cachingCloser) Read(p []byte) (n int, err error) {
	return c.Read(p)
}

func (c cachingCloser) Close() error {
	go memcache.Set(c.ctx, &memcache.Item{
		Key:   c.key,
		Value: c.buf.Bytes(),
	})
	return c.rc.Close()
}

func fromCache(ctx context.Context, bucket, name string) (io.ReadCloser, error) {
	item, err := memcache.Get(ctx, key(bucket, name))
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bytes.NewReader(item.Value)), nil
}

func fromStorage(ctx context.Context, bucket, name string) (io.ReadCloser, error) {
	rc, err := storage.NewReader(CloudContext(ctx), bucket, name)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	var crc = struct {
		io.Reader
		io.Closer
	}{
		io.TeeReader(rc, buf),
		cachingCloser{ctx: ctx, rc: rc, buf: buf, key: key(bucket, name)},
	}
	return crc, nil
}

func key(bucket, name string) string {
	return "gcs://" + bucket + "/" + name
}

func Expose(bucket, name string, expiry time.Time) (string, error) {
	if jwtc == nil {
		return "", errors.New("JWT configuration invalid")
	}
	opts := &storage.SignedURLOptions{
		GoogleAccessID: jwtc.Email,
		PrivateKey:     jwtc.PrivateKey,
		Method:         "GET",
		Expires:        expiry,
	}
	return storage.SignedURL(bucket, name, opts)
}
