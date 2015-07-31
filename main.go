package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coduno/app/controllers"
	"github.com/coduno/app/subscription"
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

func main() {
	http.HandleFunc("/status", controllers.Status)
	http.HandleFunc("/_ah/mail/", controllers.ReceiveMail)

	r := mux.NewRouter()
	r.HandleFunc("/subscriptions", secure(subscription.Subscriptions))

	r.HandleFunc("/api/code/download", setup(controllers.Template))
	r.HandleFunc("/api/invitations", setup(controllers.Invitation))

	r.HandleFunc("/api/challenges/{key}", setup(controllers.ChallengeByKey))

	r.HandleFunc("/api/companies", setup(controllers.PostCompany))
	r.HandleFunc("/api/companies/{key}/challenges", setup(controllers.GetChallengesForCompany))

	r.HandleFunc("/api/task/{key}", setup(controllers.TaskByKey))

	r.HandleFunc("/api/results", setup(controllers.CreateResult))
	r.HandleFunc("/api/results/{key}/submissions", setup(controllers.PostSubmission))

	if appengine.IsDevAppServer() {
		r.HandleFunc("/api/mock", controllers.Mock)
	}

	http.Handle("/", r)
	appengine.Main()
}

// secure is a basic wrapper that is extremely general and takes care of baseline
// features, such as tightly timed HSTS for all requests and automatic upgrades from
// HTTP to HTTPS. All outbound flows SHOULD be wrapped.
func secure(h http.HandlerFunc) http.HandlerFunc {
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

		// Redirect all HTTP requests to their HTTPS version.
		// This uses a permanent redirect to make clients adjust their bookmarks.
		if r.URL.Scheme != "https" {
			version := appengine.VersionID(appengine.NewContext(r))
			version = version[0:strings.Index(version, ".")]

			// Using www here, because cod.uno will redirect to www
			// anyway.
			host := "www.cod.uno"
			if version != "master" {
				host = version + "-dot-coduno.appspot.com"
			}

			location := "https://" + host + r.URL.Path
			http.Redirect(w, r, location, http.StatusMovedPermanently)
			return
		}

		h(w, r)
	})
}

// Rudimentary CORS checking. See
// https://developer.mozilla.org/docs/Web/HTTP/Access_control_CORS
func cors(h http.HandlerFunc) http.HandlerFunc {
	if appengine.IsDevAppServer() {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "" {
			h(w, r)
			return
		}

		// only allow CORS on localhost for development
		if !strings.HasPrefix(origin, "http://localhost") {
			// TODO(flowlo): We are probably not answering this request.
			// Is that acceptable? How to answer CORS correctly in case
			// we do not want to accept?
			return
		}

		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
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
	return secure(cors(guard(auth(h))))
}
