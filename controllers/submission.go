package controllers

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
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

func init() {
	router.Handle("/results/{resultKey}/tasks/{taskKey}/submissions", ContextHandlerFunc(PostSubmission))
	router.Handle("/results/{resultKey}/finalSubmissions/{index}", ContextHandlerFunc(FinalSubmission))
	router.Handle("/results/{resultKey}/task/{taskKey}/submissions", ContextHandlerFunc(GetSubmissionsForResult))
	router.Handle("/submissions/{key}", ContextHandlerFunc(GetSubmissionByKey))
	router.Handle("/submissions/{key}/testresults", ContextHandlerFunc(GetTestResultsForSubmission))
}

// PostSubmission creates a new submission.
func PostSubmission(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, nil
	}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return http.StatusBadRequest, err
	}
	if !strings.HasPrefix(mediaType, "multipart/") {
		return http.StatusUnsupportedMediaType, nil
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

	if err := r.ParseMultipartForm(16 << 20); err != nil {
		return http.StatusBadRequest, err
	}

	files, ok := r.MultipartForm.File["files"]
	if !ok {
		return http.StatusBadRequest, errors.New("missing files")
	}

	var task model.Task
	if err = datastore.Get(ctx, taskKey, &task); err != nil {
		return http.StatusNotFound, err
	}

	// Furthermore, the name of the GCS object is derived from the of the
	// encapsulating Submission. To avoid race conditions, allocate an ID.
	low, _, err := datastore.AllocateIDs(ctx, model.SubmissionKind, resultKey, 1)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	submissionKey := datastore.NewKey(ctx, model.SubmissionKind, "", low, resultKey)
	storedCode := model.StoredObject{
		Bucket: util.SubmissionBucket(),
		Name:   nameObject(submissionKey) + "/Code/",
	}
	submission := model.Submission{
		Task:     taskKey,
		Time:     time.Now(),
		Language: detectLanguage(files),
		Code:     storedCode,
	}

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

	prrs, pwrs := multiPipe(len(tests))

	go maketar(pwrs, files)

	for i, t := range tests {
		go func(i int, t model.Test) {
			if err := test.Tester(t.Tester).Call(ctx, *t.Key(testKeys[i]), *submission.Key(submissionKey), prrs[i]); err != nil {
				log.Warningf(ctx, "%s", err)
			}
		}(i, t)
	}

	if err := upload(util.CloudContext(ctx), storedCode.Bucket, storedCode.Name, files); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func GetSubmissionsForResult(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	// TODO(victorbalan): Check if user is company or if user is parent of result
	// else return http.StatusUnauthorized
	_, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	resultKey, err := datastore.DecodeKey(mux.Vars(r)["resultKey"])
	if err != nil {
		return http.StatusNotFound, err
	}

	taskKey, err := datastore.DecodeKey(mux.Vars(r)["taskKey"])
	if err != nil {
		return http.StatusNotFound, err
	}

	var submissions model.Submissions
	var keys []*datastore.Key
	keys, err = model.NewQueryForSubmission().
		Ancestor(resultKey).
		Filter("Task =", taskKey).
		Order("Time").
		GetAll(ctx, &submissions)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(submissions.Key(keys))
	return http.StatusOK, nil
}

func GetSubmissionByKey(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	// TODO(victorbalan): Check if user is company or if user is parent of result
	// else return http.StatusUnauthorized
	_, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	submissionKey, err := datastore.DecodeKey(mux.Vars(r)["key"])
	if err != nil {
		return http.StatusNotFound, err
	}

	var submission model.Submission
	if err = datastore.Get(ctx, submissionKey, &submission); err != nil {
		return
	}

	codeFilesURLs, err := util.ExposeMultiURL(ctx, submission.Code.Bucket, submission.Code.Name)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	s := struct {
		Time     time.Time
		CodeURLs []string
		Language string
	}{
		Time:     submission.Time,
		Language: submission.Language,
		CodeURLs: codeFilesURLs,
	}

	json.NewEncoder(w).Encode(s)
	return http.StatusOK, nil
}

func GetTestResultsForSubmission(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	// TODO(victorbalan): Check if user is company or if user is parent of result
	// else return http.StatusUnauthorized
	_, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	submissionKey, err := datastore.DecodeKey(mux.Vars(r)["key"])
	if err != nil {
		return http.StatusNotFound, err
	}

	keys, err := datastore.NewQuery("").
		Ancestor(submissionKey).
		Filter("__key__ >", submissionKey).
		KeysOnly().
		GetAll(ctx, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if len(keys) == 0 {
		json.NewEncoder(w).Encode([]string{})
		return http.StatusOK, nil
	}

	switch keys[0].Kind() {
	case model.JunitTestResultKind:
		var results model.JunitTestResults
		_, err = datastore.NewQuery(keys[0].Kind()).
			Ancestor(submissionKey).
			GetAll(ctx, &results)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		json.NewEncoder(w).Encode(results)
	case model.DiffTestResultKind:
		var results model.DiffTestResults
		_, err = datastore.NewQuery(keys[0].Kind()).
			Ancestor(submissionKey).
			GetAll(ctx, &results)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		json.NewEncoder(w).Encode(results)
	default:
		w.Write([]byte("[]"))
	}
	return http.StatusOK, nil
}

type multiError []error

func (e multiError) Error() string {
	s := make([]string, 0, len(e))
	for i := range e {
		s = append(s, e[i].Error())
	}
	return strings.Join(s, "\n")
}

func upload(ctx context.Context, bucket, base string, files []*multipart.FileHeader) multiError {
	errc := make(chan error)
	var errs []error

	for _, fh := range files {
		go func(fh *multipart.FileHeader) {
			f, err := fh.Open()
			if err != nil {
				errc <- err
				return
			}

			name := base + fh.Filename
			wc := storage.NewWriter(ctx, bucket, name)

			wc.ObjectAttrs = defaultObjectAttrs(path.Base(fh.Filename))

			if _, err := io.Copy(wc, f); err != nil {
				errc <- err
				return
			}

			if err := wc.Close(); err != nil {
				errc <- err
				return
			}

			if err := f.Close(); err != nil {
				errc <- err
				return
			}

			errc <- nil
		}(fh)
	}

	for range files {
		err := <-errc
		if err != nil {
			errs = append(errs, err)
		}
	}

	close(errc)
	return errs
}

func maketar(pw []*io.PipeWriter, files []*multipart.FileHeader) error {
	var w []io.Writer
	for i := range pw {
		w = append(w, pw[i])
		defer pw[i].Close()
	}
	tarw := tar.NewWriter(io.MultiWriter(w...))
	defer tarw.Close()

	sizeFunc := func(s io.Seeker) int64 {
		size, err := s.Seek(0, os.SEEK_END)
		if err != nil {
			return -1
		}
		if _, err = s.Seek(0, os.SEEK_SET); err != nil {
			return -1
		}
		return size
	}

	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			return err
		}
		size := sizeFunc(f)
		if size < 0 {
			return errors.New("seeker can't seek")
		}
		tarw.WriteHeader(&tar.Header{
			Name: fh.Filename,
			Mode: 0600,
			Size: size,
		})
		io.Copy(tarw, f)
		f.Close()
	}
	return nil
}

func multiPipe(n int) ([]*io.PipeReader, []*io.PipeWriter) {
	wrs := make([]*io.PipeWriter, n)
	rrs := make([]*io.PipeReader, n)
	for i := 0; i < n; i++ {
		rrs[i], wrs[i] = io.Pipe()
	}
	return rrs, wrs
}

func detectLanguage(files []*multipart.FileHeader) string {
	l := ""
	m := map[string]int{
		"py":   0,
		"java": 0,
		"c":    0,
		"cpp":  0,
	}
	max := 0

	for _, fh := range files {
		name := fh.Filename
		dot := strings.LastIndex(name, ".") + 1
		if dot == 0 || dot == len(name) {
			continue
		}
		ext := name[dot:]
		cnt, ok := m[ext]
		if !ok {
			continue
		}
		cnt++
		if cnt > max {
			l = ext
		}
		m[ext] = cnt
	}
	return l
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
	// this will give you a major headache when working with
	// gcsfuse. Just don't.
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
