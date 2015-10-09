package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/coduno/api/db"
	"github.com/coduno/api/util/passenger"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

// Holds routes to all controllers' handlers.
var router = mux.NewRouter()

func Handler() http.Handler {
	return router
}

func decode(r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&dst); err != nil {
		return err
	}
	return nil
}

type SimpleContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (h SimpleContextHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if p := recover(); p != nil {
			deal(nil, w, r, http.StatusInternalServerError, errors.New(fmt.Sprint(p)))
		}
	}()

	ctx := appengine.NewContext(r)
	hsts(w)
	if !cors(w, r) {
		return
	}

	// Add authentication metadata.
	ctx, err := passenger.NewContextFromRequest(ctx, r)
	if err != nil {
		deal(ctx, w, r, http.StatusInternalServerError, err)
		return
	}

	h(ctx, w, r)
}

// ContextHandlerFunc is similar to a http.HandlerFunc, but also gets passed
// the current context.
// To ease error handling, a ContextHandleFunc must return a HTTP status
// code and an error. Still, the handler is allowed to write a response.
type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request) (int, error)

func (h ContextHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if p := recover(); p != nil {
			deal(nil, w, r, http.StatusInternalServerError, errors.New(fmt.Sprint(p)))
		}
	}()

	ctx := appengine.NewContext(r)
	hsts(w)
	if !cors(w, r) {
		return
	}

	// Add authentication metadata.
	ctx, err := passenger.NewContextFromRequest(ctx, r)

	buf := bufferedResponseWriter{
		b: new(bytes.Buffer),
		w: w,
		s: 0,
	}

	status, err := h(ctx, buf, r)

	if status == 0 {
		status = buf.s
	}

	// No error and a low status code means
	// everything went well.
	if err == nil && status < 400 {
		if status == 0 {
			status = http.StatusOK
		}
		buf.flush(status)
		return
	}

	deal(ctx, w, r, status, err)
}

func hsts(w http.ResponseWriter) {
	if !appengine.IsDevAppServer() {
		// Protect against HTTP downgrade attacks by explicitly telling
		// clients to use HTTPS.
		// max-age is computed to match the expiration date of our TLS
		// certificate.
		// https://developer.mozilla.org/docs/Web/Security/HTTP_strict_transport_security
		// This is only set on production.
		invalidity := time.Date(2017, time.July, 15, 8, 30, 21, 0, time.UTC)
		maxAge := invalidity.Sub(time.Now()).Seconds()
		w.Header().Set("Strict-Transport-Security", fmt.Sprintf("max-age=%d", int(maxAge)))
	}
}

// Rudimentary CORS checking. See
// https://developer.mozilla.org/docs/Web/HTTP/Access_control_CORS
func cors(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")

	// If the client has not provided it's origin, the
	// request will be answered in any case.
	if origin == "" {
		return true
	}

	// Only allow our own origin if not on development server.
	if !appengine.IsDevAppServer() && origin != "https://app.cod.uno" {
		http.Error(w, "Invalid Origin", http.StatusUnauthorized)
		return false
	}

	// We have a nice CORS established, so set appropriate headers.
	// TODO(flowlo): Figure out how to send correct methods.
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", origin)

	// If this is an OPTIONS request, we answer it
	// immediately and do not bother higher level handlers.
	if r.Method == "OPTIONS" {
		w.Write([]byte("OK"))
		return false
	}

	return true
}

type bufferedResponseWriter struct {
	b *bytes.Buffer
	w http.ResponseWriter
	s int
}

func (w bufferedResponseWriter) flush(status int) {
	w.w.WriteHeader(status)
	io.Copy(w.w, w.b)
}

func (w bufferedResponseWriter) WriteHeader(status int) {
	w.s = status
}

func (w bufferedResponseWriter) Write(p []byte) (n int, err error) {
	return w.b.Write(p)
}

func (w bufferedResponseWriter) Header() http.Header {
	return w.w.Header()
}

type trace struct {
	e error
	t []byte
}

func (t trace) Error() string {
	return fmt.Sprintf("%s\n%s", t.e, t.t)
}

// tracable wraps the passed error and generates a new error that will
// expand into a full stack trace. Be aware that this is expensive,
// as it will stop the world to collect the trace, and should be
// used with caution!
func tracable(err error) error {
	r := trace{
		e: err,
		t: make([]byte, 4096),
	}
	runtime.Stack(r.t, false)
	return r
}

func respond(ctx context.Context, w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	if err, ok := data.(error); ok || status < 200 || status >= 300 {
		deal(ctx, w, r, status, err)
		return
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	io.Copy(w, &buf)
}

// deal makes the response in error cases somewhat nicer. It will try
// to figure out what actually went wrong and inform the user.
// It should not be called if the request went fine. If status is below
// 400, and err is not nil, it will assume an internal server error.
// Generally, if you pass a nil error, don't expect deal to do anything
// useful.
func deal(ctx context.Context, w http.ResponseWriter, r *http.Request, status int, err error) {
	// Getting an error and a status code blow 400 is somewhat paradox.
	// Also, if the status is the zero value, assume that we're dealing
	// with an internal server error.
	if err != nil && status < 400 || status == 0 {
		switch err {
		case db.NotValid:
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}
	}

	codunoError := struct {
		Message,
		Reason,
		RequestID,
		StatusText,
		Trace string
		Status int
	}{}

	codunoError.Status = status
	codunoError.StatusText = http.StatusText(status)
	if ctx != nil {
		codunoError.RequestID = appengine.RequestID(ctx)
	}

	w.WriteHeader(status)

	// If we don't have an error it's really hard to make sense.
	if err == nil {
		json.NewEncoder(w).Encode(codunoError)
		return
	}

	if appengine.IsOverQuota(err) {
		codunoError.Reason = "Over Quota"
	} else if appengine.IsTimeoutError(err) {
		codunoError.Reason = "Timeout Error"
	} else {
		codunoError.Reason = err.Error()
	}

	if t, ok := err.(trace); ok {
		codunoError.Trace = strings.Replace(string(t.t), "\n", "\n\t", -1)
	}

	json.NewEncoder(w).Encode(codunoError)
}
