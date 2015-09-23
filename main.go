package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coduno/api/controllers"
	"github.com/coduno/api/util/passenger"
	"github.com/coduno/api/ws"
	"github.com/gorilla/mux"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/cloud/compute/metadata"
)

// ContextHandlerFunc is similar to a http.HandlerFunc, but also gets passed
// the current context.
// To ease error handling, a ContextHandleFunc must return a HTTP status
// code and an error. Still, the handler is allowed to write a response.
type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request) (int, error)

func main() {
	go http.ListenAndServe(":8090", http.HandlerFunc(ws.Handle))

	http.HandleFunc("/_ah/mail/", controllers.ReceiveMail)
	http.HandleFunc("/cert", hsts(controllers.Certificate))
	http.HandleFunc("/status", hsts(controllers.Status))
	http.HandleFunc("/ip", hsts(ip))

	r := mux.NewRouter()
	r.HandleFunc("/subscriptions", hsts(controllers.Subscriptions))

	r.HandleFunc("/code/download/{name}/{language}", setup(controllers.Template))
	r.HandleFunc("/invitations", setup(controllers.Invitation))

	r.HandleFunc("/tokens", setup(controllers.Tokens))
	r.HandleFunc("/tokens/collect", setup(controllers.CollectTokens))

	r.HandleFunc("/challenges", setup(controllers.CreateChallenge))
	r.HandleFunc("/challenges/{key}", setup(controllers.ChallengeByKey))
	r.HandleFunc("/challenges/{key}/results", setup(controllers.GetResultsByChallenge))

	r.HandleFunc("/companies", setup(controllers.PostCompany))
	r.HandleFunc("/companies/{key}/challenges", setup(controllers.GetChallengesForCompany))
	r.HandleFunc("/companies/{key}/users", setup(controllers.GetUsersByCompany))

	r.HandleFunc("/mock", controllers.Mock)

	r.HandleFunc("/profiles/{key}", setup(controllers.GetProfileByKey))
	r.HandleFunc("/profiles/{key}", setup(controllers.DeleteProfile))
	r.HandleFunc("/profiles/{key}/challenges", setup(controllers.GetChallengesForProfile))

	r.HandleFunc("/results", setup(controllers.CreateResult))
	r.HandleFunc("/results/{resultKey}/tasks/{taskKey}/submissions", setup(controllers.PostSubmission))
	r.HandleFunc("/results/{resultKey}/finalSubmissions/{index}", setup(controllers.FinalSubmission))
	r.HandleFunc("/results/{resultKey}", setup(controllers.GetResult))
	r.HandleFunc("/results/user/{userKey}/challenge/{challengeKey}", setup(controllers.GetResultForUserChallenge))

	r.HandleFunc("/user/company", setup(controllers.GetCompanyByUser))
	r.HandleFunc("/user", setup(controllers.WhoAmI))
	r.HandleFunc("/users", setup(controllers.User))
	r.HandleFunc("/users/{key}", setup(controllers.GetUser))
	r.HandleFunc("/users/{key}/profile", setup(controllers.GetProfileForUser))

	r.HandleFunc("/tasks/{key}", setup(controllers.TaskByKey))
	r.HandleFunc("/tasks/{key}/tests", setup(controllers.TestsByTaskKey))
	r.HandleFunc("/tasks", setup(controllers.Tasks))

	r.HandleFunc("/whoami", setup(controllers.WhoAmI))

	http.Handle("/", r)
	appengine.Main()
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

func ip(w http.ResponseWriter, r *http.Request) {
	ip := "127.0.0.1"
	if metadata.OnGCE() {
		var err error
		ip, err = metadata.InternalIP()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	w.Write([]byte(ip))
}
