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

// ContextHandleFunc is similar to a HandleFunc, but also gets passed
// the current context.
type ContextHandleFunc func(context.Context, http.ResponseWriter, *http.Request)

// HandleFunc describes a function that can handle HTTP requests and respond
// to them correctly.
type HandleFunc func(http.ResponseWriter, *http.Request)

func main() {
	http.HandleFunc("/status", controllers.Status)
	http.HandleFunc("/_ah/mail/", controllers.ReceiveMail)

	r := mux.NewRouter()
	r.HandleFunc("/api/companies/{id}/challenges", setup(controllers.GetChallengesForCompany))
	r.HandleFunc("/api/challenges/{id}", setup(controllers.GetChallengeByID))
	r.HandleFunc("/api/results", setup(controllers.CreateResult))
	r.HandleFunc("/api/code/download", setup(controllers.DownloadTemplate))
	r.HandleFunc("/api/companies", setup(controllers.CreateCompany))
	r.HandleFunc("/api/company/login", setup(controllers.CompanyLogin))
	r.HandleFunc("/api/engineurl", secure(controllers.Engine))
	r.HandleFunc("/api/fingerprints", setup(controllers.HandleFingerprints))
	r.HandleFunc("/api/invitations", setup(controllers.Invitation))
	r.HandleFunc("/api/mock", controllers.MockData)
	r.HandleFunc("/api/mockCompany", controllers.MockCompany)
	r.HandleFunc("/subscriptions", secure(subscription.Subscriptions))
	r.HandleFunc("/api/task/{id}", setup(controllers.GetTaskByKey))
	r.HandleFunc("/api/results/{id}/submission", setup(controllers.PostSubmission))
	http.Handle("/", r)
	appengine.Main()
}

// setupBaseHandler is a basic wrapper that is extremely general and takes care of baseline
// features, such as tightly timed HSTS for all requests and automatic upgrades from
// HTTP to HTTPS. All outbound flows SHOULD be wrapped.
func secure(h HandleFunc) HandleFunc {
	if appengine.IsDevAppServer() {
		return h
	}

	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

// Rudimentary CORS checking. See
// https://developer.mozilla.org/docs/Web/HTTP/Access_control_CORS
func cors(h HandleFunc) HandleFunc {
	if appengine.IsDevAppServer() {
		return h
	}
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

// auth is there to associate a user with the incoming request.
func auth(h ContextHandleFunc) HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := passenger.NewContextFromRequest(appengine.NewContext(r), r)
		if err != nil {
			log.Debugf(ctx, "auth: "+err.Error())
		}
		h(ctx, w, r)
	}
}

// setup is the default wrapper for any HandleFunc that talks to
// the outside. It will wrap h in scure, cors and auth.
func setup(h ContextHandleFunc) HandleFunc {
	return secure(cors(auth(h)))
}
