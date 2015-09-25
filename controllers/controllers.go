package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coduno/engine/passenger"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// ContextHandlerFunc is similar to a http.HandlerFunc, but also gets passed
// the current context.
// To ease error handling, a ContextHandleFunc must return a HTTP status
// code and an error. Still, the handler is allowed to write a response.
type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request) (int, error)

var router = mux.NewRouter()

func Handler() http.Handler {
	return router
}

// hsts is a basic wrapper that takes care of tightly timed HSTS for all requests.
// All outbound flows should be wrapped.
func hsts(h http.HandlerFunc) http.HandlerFunc {
	if appengine.IsDevAppServer() {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Protect against HTTP downgrade attacks by explicitly telling
		// clients to use HTTPS.
		// max-age is computed to match the expiration date of our TLS
		// certificate.
		// https://developer.mozilla.org/docs/Web/Security/HTTP_strict_transport_security
		// This is only set on production.
		invalidity := time.Date(2017, time.July, 15, 8, 30, 21, 0, time.UTC)
		maxAge := invalidity.Sub(time.Now()).Seconds()
		w.Header().Set("Strict-Transport-Security", fmt.Sprintf("max-age=%d", int(maxAge)))
		h(w, r)
	})
}

// Rudimentary CORS checking. See
// https://developer.mozilla.org/docs/Web/HTTP/Access_control_CORS
func cors(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if !appengine.IsDevAppServer() {
			if origin == "" {
				h(w, r)
				return
			}

			if !strings.HasPrefix(origin, "https://app.cod.uno") {
				http.Error(w, "Invalid CORS request", http.StatusUnauthorized)
				return
			}
		}
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", origin)

		if r.Method == "OPTIONS" {
			w.Write([]byte("OK"))
			return
		}

		h(w, r)
	})
}

// auth is there to associate a user with the incoming request.
func auth(h ContextHandlerFunc) ContextHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
		ctx, err = passenger.NewContextFromRequest(ctx, r)
		if err != nil {
			log.Warningf(ctx, "auth: "+err.Error())
		}
		return h(ctx, w, r)
	}
}

func guard(h ContextHandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status, err := h(appengine.NewContext(r), w, r)

		if err != nil {
			http.Error(w, err.Error(), status)
		} else if status >= 400 {
			http.Error(w, http.StatusText(status), status)
		}
	})
}

// setup is the default wrapper for any ContextHandlerFunc that talks to
// the outside. It will wrap h in scure, cors and auth.
func setup(h ContextHandlerFunc) http.HandlerFunc {
	return hsts(cors(guard(auth(h))))
}
