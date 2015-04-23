package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"mail"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"github.com/coduno/app/controllers"
	"github.com/coduno/app/util"
)

var gitlabToken = "YHQiqMx3qUfj8_FxpFe4"

type Handler func(http.ResponseWriter, *http.Request)

// A basic wrapper that is extremely general and takes care of baseline features, such
// as tightly timed HSTS for all requests and automatic upgrades from HTTP to HTTPS.
// All outbound flows SHOULD be wrapped.
func setupHandler(handler func(http.ResponseWriter, *http.Request, context.Context)) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// Redirect all HTTP requests to their HTTPS version.
		// This uses a permanent redirect to make clients adjust their bookmarks.
		if r.URL.Scheme != "https" {
			location := r.URL
			location.Scheme = "https"
			http.Redirect(w, r, location.String(), http.StatusMovedPermanently)
			return
		}

		// Protect against HTTP downgrade attacks by explicitly telling
		// clients to use HTTPS.
		// max-age is computed to match the expiration date of our TLS
		// certificate (minus approx. one day buffer).
		// https://developer.mozilla.org/docs/Web/Security/HTTP_strict_transport_security
		invalidity := time.Date(2016, time.January, 3, 0, 59, 59, 0, time.UTC)
		maxAge := invalidity.Sub(time.Now()).Seconds()
		w.Header().Set("Strict-Transport-Security", fmt.Sprintf("max-age=%d", int(maxAge)))

		// appengine.NewContext is cheap, and context is needed in many
		// handlers, so create one here.
		c := appengine.NewContext(r)

		// All wrapping is done, call the original handler.
		handler(w, r, c)
	}
}

func main() {
	http.HandleFunc("/api/token", setupHandler(token))
	http.HandleFunc("/_ah/health", health)
	http.HandleFunc("/_ah/start", start)
	http.HandleFunc("/_ah/stop", stop)
	http.ListenAndServe(":8080", nil)
	http.HandleFunc("/push", setupHandler(controllers.Push))
	http.HandleFunc("/subscriptions", setupHandler(mail.Subscriptions))
}

func health(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "OK")
}

func start(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "OK")
}

func stop(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Stopping...")
}

func generateToken() (string, error) {
	token := make([]byte, 64)

	if _, err := rand.Read(token); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(token), nil
}

func token(w http.ResponseWriter, req *http.Request, c context.Context) {
	if req.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Invalid method.", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(body) < 1 {
		http.Error(w, "Invalid body.", http.StatusBadRequest)
	}

	if err, username := authenticate(req); err == nil {
		token, err := generateToken()
		if err != nil {
			http.Error(w, "Failed to generate token.", http.StatusInternalServerError)
		} else {
			query := url.Values{}
			query.Add("id", "2")
			query.Add("title", username)
			query.Add("key", string(body))

			gitlabReq, _ := http.NewRequest("POST", "http://git.cod.uno/api/v3/users/2/keys", strings.NewReader(query.Encode()))

			gitlabReq.Header = map[string][]string{
				"PRIVATE-TOKEN": {gitlabToken},
				"Content-Type":  {"application/x-www-form-urlencoded"},
				"Accept":        {"application/json"},
			}

			client := urlfetch.Client(c)
			res, err := client.Do(gitlabReq)
			if err != nil {
				log.Debugf(c, err.Error())
			} else {
				result, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Debugf(c, err.Error())
				} else {
					log.Debugf(c, string(result))
				}
			}
			defer res.Body.Close()

			log.Infof(c, "Generated token for '%s'", username)
			fmt.Fprintf(w, token)
		}
	} else {
		// This could be either invalid/missing credentials or
		// the database failing, so let's issue a warning.
		log.Warningf(c, err.Error())
		http.Error(w, "Invalid credentials.", http.StatusForbidden)
	}
}

func authenticate(req *http.Request) (error, string) {
	username, _, ok := util.BasicAuth(req)
	if !ok {
		return errors.New("No Authorization header present."), ""
	}
	return errors.New("No authorization backend present."), username
}
