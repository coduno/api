package controllers

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/cloud/storage"

	"github.com/coduno/api/model"
	"github.com/coduno/api/test"
	"github.com/coduno/api/util"
	"github.com/coduno/api/util/passenger"
	"github.com/gorilla/mux"

	"golang.org/x/net/context"
)

func init() {
	router.Handle("/results/{resultKey}/tasks/{taskKey}/submissions", ContextHandlerFunc(PostSubmission))
	router.Handle("/results/{resultKey}/finalSubmissions/{index}", ContextHandlerFunc(FinalSubmission))
}

// PostSubmission creates a new submission.
func PostSubmission(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, nil
	}

	var body = struct {
		Code     string
		Language string
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return http.StatusBadRequest, err
	}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	resultKey, err := datastore.DecodeKey(mux.Vars(r)["resultKey"])
	if err != nil {
		return http.StatusNotFound, err
	}

	if !util.HasParent(p.User, resultKey) {
		return http.StatusBadRequest, errors.New("cannot submit answer for other users")
	}

	taskKey, err := datastore.DecodeKey(mux.Vars(r)["taskKey"])
	if err != nil {
		return http.StatusNotFound, err
	}

	var task model.Task
	if err = datastore.Get(ctx, taskKey, &task); err != nil {
		return http.StatusInternalServerError, err
	}

	// Furthermore, the name of the GCS object is derived from the of the
	// encapsulating Submission. To avoid race conditions, allocate an ID.
	low, _, err := datastore.AllocateIDs(ctx, model.SubmissionKind, resultKey, 1)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	submissionKey := datastore.NewKey(ctx, model.SubmissionKind, "", low, resultKey)

	submission := model.Submission{
		Task: taskKey,
		Time: time.Now(),
	}

	var ball io.Reader

	if body.Code != "" {
		submission.Code, ball, err = ingest(ctx, submissionKey, body.Code, body.Language)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		submission.Language = body.Language
	}

	// Set the submission in stone.
	if _, err = datastore.Put(ctx, submissionKey, &submission); err != nil {
		return http.StatusInternalServerError, err
	}

	var tests model.Tests
	testKeys, err := model.NewQueryForTest().
		Ancestor(taskKey).
		GetAll(ctx, &tests)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	for i, t := range tests {
		go func(i int, t model.Test) {
			// TODO(victorbalan, flowlo): Error handling
			if err := test.Tester(t.Tester).Call(ctx, *t.Key(testKeys[i]), *submission.Key(submissionKey), ball); err != nil {
				log.Warningf(ctx, "%s", err)
			}
		}(i, t)
	}

	// TODO(flowlo): Return something meaningful.

	return http.StatusOK, nil
}

// ingest will take an upload, pipe it to GCS and return a tarball.
func ingest(ctx context.Context, key *datastore.Key, code, language string) (model.StoredObject, io.Reader, error) {
	var o model.StoredObject

	pr, pw := io.Pipe()

	// We'll be storing code, so check what
	// name the file should have in GCS and TAR.
	fn, ok := util.FileNames[language]
	if !ok {
		return o, nil, errors.New("language unknown")
	}

	var gcsw io.WriteCloser

	if appengine.IsDevAppServer() {
		// All hail /dev/null!
		gcsw = struct {
			io.Writer
			io.Closer
		}{
			ioutil.Discard,
			ioutil.NopCloser(nil),
		}
	} else {
		// Let's be serious and upload the code to GCS.
		// TODO(flowlo): Limit this writer, or limit the uploaded code
		// at some previous point.
		o = model.StoredObject{
			Bucket: util.SubmissionBucket(),
			Name:   nameObject(key) + "/Code/" + fn,
		}

		gcsw = storage.NewWriter(util.CloudContext(ctx), o.Bucket, o.Name)
		if gcsw == nil {
			return o, nil, errors.New("cannot obtain writer to gcs")
		}

		gcsw.(*storage.Writer).ObjectAttrs = defaultObjectAttrs(fn)
	}

	bc := []byte(code)

	// Create a TAR stream.
	tarw := tar.NewWriter(pw)

	// Stream to TAR and GCS at the same time.
	// TODO(flowlo): Pass code as io.Reader.
	go func() {
		tarw.WriteHeader(&tar.Header{
			Name: fn,
			Mode: 0600,
			Size: int64(len(bc)),
		})
		io.Copy(io.MultiWriter(tarw, gcsw), bytes.NewReader(bc))
		gcsw.Close()
		tarw.Close()
		pw.Close()
	}()

	return o, pr, nil
}

// FinalSubmission makes the last submission final.
func FinalSubmission(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, nil
	}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	var resultKey *datastore.Key
	if resultKey, err = datastore.DecodeKey(mux.Vars(r)["resultKey"]); err != nil {
		return http.StatusInternalServerError, err
	}

	if !util.HasParent(p.User, resultKey) {
		return http.StatusBadRequest, errors.New("cannot submit answer for other users")
	}

	var index int
	if index, err = strconv.Atoi(mux.Vars(r)["index"]); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(r.URL.Query()["submissionKey"]) == 0 {
		return http.StatusOK, nil
	}
	var submissionKey *datastore.Key
	if submissionKey, err = datastore.DecodeKey(r.URL.Query()["submissionKey"][0]); err != nil {
		return http.StatusInternalServerError, err
	}

	var result model.Result
	if err = datastore.Get(ctx, resultKey, &result); err != nil {
		return http.StatusInternalServerError, err
	}

	result.FinalSubmissions[index] = submissionKey

	if _, err = result.Put(ctx, resultKey); err != nil {
		return http.StatusInternalServerError, err
	}
	w.Write([]byte("OK"))
	return
}

func nameObject(key *datastore.Key) string {
	name := ""
	for key != nil {
		id := key.StringID()
		if id == "" {
			id = strconv.FormatInt(key.IntID(), 10)
		}
		name = "/" + key.Kind() + "/" + id + name
		key = key.Parent()
	}
	// NOTE: The name of a GCS object must not be prefixed "/",
	// this will give you a major headache.
	return name[1:]
}

func defaultObjectAttrs(disposition string) storage.ObjectAttrs {
	// TODO(flowlo): Establish ACLs?
	return storage.ObjectAttrs{
		ContentType:        "text/plain", // TODO(flowlo): Content types per language?
		ContentLanguage:    "",           // TODO(flowlo): Does it make sense to set this?
		ContentEncoding:    "utf-8",
		CacheControl:       "max-age=31556926", // Aggressive caching for one year
		ContentDisposition: "attachment; filename=\"" + disposition + "\"",
	}
}
