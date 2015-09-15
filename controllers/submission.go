package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

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

	if body.Code != "" {
		submission.Code, err = store(ctx, submissionKey, body.Code, body.Language)
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
	_, err = model.NewQueryForTest().
		Ancestor(taskKey).
		GetAll(ctx, &tests)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	for _, t := range tests {
		if err := test.Tester(t.Tester).Call(ctx, t.Params, *submission.Key(submissionKey)); err != nil {
			log.Warningf(ctx, "%s", err)
			continue
		}
	}

	// TODO(flowlo): Return something meaningful.

	return http.StatusOK, nil
}

func store(ctx context.Context, key *datastore.Key, code, language string) (model.StoredObject, error) {
	o := model.StoredObject{}

	// We'll be storing code, so check what name the file should have in GCS.
	fn, ok := util.FileNames[language]
	if !ok {
		return o, errors.New("language unknown")
	}

	submissionBucket := util.SubmissionBucket()

	// Now, construct the object.
	o = model.StoredObject{
		Bucket: submissionBucket,
		Name:   nameObject(key) + "/Code/" + fn,
	}

	// Upload the code to GCS.
	// TODO(flowlo): Limit this writer, or limit the uploaded code
	// at some previous point.
	gcs := storage.NewWriter(util.CloudContext(ctx), o.Bucket, o.Name)
	if gcs == nil {
		return o, errors.New("cannot obtain writer to gcs")
	}
	gcs.ObjectAttrs = defaultObjectAttrs(fn)

	if _, err := io.WriteString(gcs, code); err != nil {
		return o, err
	}

	return o, gcs.Close()
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
