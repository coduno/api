package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coduno/app/controllers"
	"github.com/coduno/app/mail"
	"github.com/coduno/app/models"
	"github.com/gorilla/mux"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

var gitlabToken = "YHQiqMx3qUfj8_FxpFe4"

// Handler is similar to a HandlerFunc, but also gets passed
// the current context.
type Handler func(http.ResponseWriter, *http.Request, context.Context)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/subscriptions", setupHandler(mail.Subscriptions))
	r.HandleFunc("/api/token", setupHandler(token))
	r.HandleFunc("/api/push", setupHandler(controllers.Push))
	r.HandleFunc("/api/code/upload", controllers.UploadCode)
	r.HandleFunc("/api/code/download", controllers.DownloadTemplate)
	r.HandleFunc("/api/token/check/{token}", setupHandler(controllers.CheckToken))
	r.HandleFunc("/api/mock", mockData)
	http.Handle("/", r)
	appengine.Main()
}

func mockData(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	company := models.Company{Name: "Catalysts"}
	companyKey, _ := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "companies", nil), &company)

	challenge := models.Challenge{Name: "Tic-Tac-Toe", Instructions: "Implenet tic tac toe input and output blah blah", Company: companyKey}
	challengeKey, _ := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "challenges", nil), &challenge)

	template := models.Template{Language: "Java", Path: "/templates/TicTacToeTemplate.java", Challenge: challengeKey}
	datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "templates", nil), &template)

	coder := models.Coder{Email: "victor.balan@cod.uno", FirstName: "Victor", LastName: "Balan"}
	coderKey, _ := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "coders", nil), &coder)

	fingerprint := models.Fingerprint{Coder: coderKey, Challenge: challengeKey, Token: "ag1812378asgvasvasdnkcasr23vva="}
	datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "fingerprints", nil), &fingerprint)
}

// setupHandler is a basic wrapper that is extremely general and takes care of baseline
// features, such as tightly timed HSTS for all requests and automatic upgrades from
// HTTP to HTTPS. All outbound flows SHOULD be wrapped.
func setupHandler(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)

		if !appengine.IsDevAppServer() {

			// Redirect all HTTP requests to their HTTPS version.
			// This uses a permanent redirect to make clients adjust their bookmarks.
			// This is only done on production.
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
			// This is only set on production.
			invalidity := time.Date(2016, time.January, 3, 0, 59, 59, 0, time.UTC)
			maxAge := invalidity.Sub(time.Now()).Seconds()
			w.Header().Set("Strict-Transport-Security", fmt.Sprintf("max-age=%d", int(maxAge)))
		}

		cors(w, r)
		handler(w, r, ctx)
	}
}

// Rudimentary CORS checking. See
// https://developer.mozilla.org/docs/Web/HTTP/Access_control_CORS
func cors(w http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")

	// only allow CORS on localhost for development
	if strings.HasPrefix(origin, "http://localhost") {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
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

	if username, err := authenticate(req); err == nil {
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

func authenticate(req *http.Request) (string, error) {
	username, _, ok := req.BasicAuth()
	if !ok {
		return "", errors.New("No Authorization header present.")
	}
	return username, errors.New("No authorization backend present.")
}
