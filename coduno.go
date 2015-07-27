package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coduno/app/controllers"
	"github.com/coduno/app/mail"
	"github.com/coduno/app/models"
	"github.com/coduno/app/status"
	"github.com/coduno/app/util/password"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

var gitlabToken = "YHQiqMx3qUfj8_FxpFe4"

const sessionID = "SESSIONID"

var cs = sessions.NewCookieStore(securecookie.GenerateRandomKey(64))

// Handler is similar to a HandlerFunc, but also gets passed
// the current context.
type Handler func(http.ResponseWriter, *http.Request, context.Context)

// HandlerWithSession is similar to Handler, but it also returns a param that
// tells if we should create a new session or not
type HandlerWithSession func(http.ResponseWriter, *http.Request, context.Context) bool

func main() {
	http.HandleFunc("/status", status.Handler)
	http.HandleFunc("/_ah/mail/", receiveMail)

	r := mux.NewRouter()
	r.HandleFunc("/api/challenges", setupHandler(controllers.GetChallengesForCompany))
	r.HandleFunc("/api/challenges/{id}", setupHandler(controllers.GetChallengeByID))
	r.HandleFunc("/api/code/download", setupHandler(controllers.DownloadTemplate))
	r.HandleFunc("/api/companies", setupHandler(controllers.CreateCompany))
	r.HandleFunc("/api/company/login", setupHandlerWithSessionStore(controllers.CompanyLogin))
	r.HandleFunc("/api/fingerprints", setupHandler(controllers.HandleFingerprints))
	r.HandleFunc("/api/token", setupHandler(token))
	r.HandleFunc("/api/token/check/{token}", setupHandlerWithSessionStore(controllers.CheckToken))
	r.HandleFunc("/api/mock", mockData)
	r.HandleFunc("/api/mockCompany", mockCompany)
	r.HandleFunc("/api/engineurl", setupHandler(getEngineURL))
	r.HandleFunc("/subscriptions", setupHandler(mail.Subscriptions))
	http.Handle("/", r)
	appengine.Main()
}

func getEngineURL(w http.ResponseWriter, req *http.Request, ctx context.Context) {
	if appengine.IsDevAppServer() {
		io.WriteString(w, "http://localhost:8081")
	} else {
		io.WriteString(w, "https://engine.cod.uno")
	}
}

func mockCompany(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	pw, _ := password.Hash([]byte("123123123123"))
	cmp := models.Company{Name: "cat", Email: "paul@cod.uno", HashedPassword: pw}
	cmp.Save(ctx)
	w.Write([]byte("Its all fine"))
}
func mockData(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	pw, _ := password.Hash([]byte("victor"))
	companyKey, _ := models.Company{Name: "Coduno", Email: "victor@cod.uno", HashedPassword: pw}.Save(ctx)

	simpleChallengeKey, _ := models.Challenge{Name: "Hello world!", Instructions: "Implement a program that outputs 'Hello World!'", Company: companyKey,
		WebInterface: "simple-console", Runner: "simple", Flags: "-image simple"}.Save(ctx)
	utChallengeKey, _ := models.Challenge{Name: "Unit testing!", Instructions: `Implement a program that has an ok() method that returns "ok",
																												an hello() method that returns "hello"
																												and an "long add(long a, long b)" method`, Company: companyKey,
		WebInterface: "javaut-console", Runner: "javaut", Flags: "-image javaut"}.Save(ctx)

	template := models.Template{Language: "Java", Path: "/templates/TicTacToeTemplate.java", Challenge: simpleChallengeKey}
	datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "templates", nil), &template)

	coderKey, _ := models.Coder{Email: "victor.balan@cod.uno", FirstName: "Victor", LastName: "Balan"}.Save(ctx)

	models.Fingerprint{Coder: coderKey, Challenge: simpleChallengeKey, Token: "deadbeefcafebabe"}.Save(ctx)

	models.Fingerprint{Coder: coderKey, Challenge: utChallengeKey, Token: "asdfbsdfbvsd13"}.Save(ctx)
}

func setupHandlerWithSessionStore(handler HandlerWithSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := setupBaseHandler(w, r)
		if handler(w, r, ctx) {
			initSession(w, r)
		}
	}
}

func setupHandler(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			cors(w, r)
			return
		}
		// session, _ := cs.Get(r, sessionID)
		// if session.IsNew {
		// 	http.Error(w, "Unauthorized session.", http.StatusUnauthorized)
		// 	return
		// }
		ctx := setupBaseHandler(w, r)
		handler(w, r, ctx)
	}
}

// setupBaseHandler is a basic wrapper that is extremely general and takes care of baseline
// features, such as tightly timed HSTS for all requests and automatic upgrades from
// HTTP to HTTPS. All outbound flows SHOULD be wrapped.
func setupBaseHandler(w http.ResponseWriter, r *http.Request) context.Context {
	ctx := appengine.NewContext(r)

	if !appengine.IsDevAppServer() {
		// Redirect all HTTP requests to their HTTPS version.
		// This uses a permanent redirect to make clients adjust their bookmarks.
		// This is only done on production.
		if r.URL.Scheme != "https" {
			version := appengine.VersionID(ctx)
			version = version[0:strings.Index(version, ".")]

			// Using www here, because cod.uno will redirect to www
			// anyway.
			host := "www.cod.uno"
			if version != "master" {
				host = version + "-dot-coduno.appspot.com"
			}

			location := "https://" + host + r.URL.Path
			http.Redirect(w, r, location, http.StatusMovedPermanently)
			return nil
		}

		// Protect against HTTP downgrade attacks by explicitly telling
		// clients to use HTTPS.
		// max-age is computed to match the expiration date of our TLS
		// certificate (minus approx. one day buffer).
		// https://developer.mozilla.org/docs/Web/Security/HTTP_strict_transport_security
		// This is only set on production.
		invalidity := time.Date(2017, time.July, 15, 8, 30, 21, 0, time.UTC)
		maxAge := invalidity.Sub(time.Now()).Seconds()
		w.Header().Set("Strict-Transport-Security", fmt.Sprintf("max-age=%d", int(maxAge)))
	}

	cors(w, r)
	return ctx
}

// initSession starts a new api session and sets the cookie header with the
// API_TOKEN required for api calls
func initSession(w http.ResponseWriter, r *http.Request) {
	session, _ := cs.New(r, sessionID)
	session.Options.MaxAge = 12 * 3600
	session.Save(r, w)
}

// Rudimentary CORS checking. See
// https://developer.mozilla.org/docs/Web/HTTP/Access_control_CORS
func cors(w http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")

	// only allow CORS on localhost for development
	if strings.HasPrefix(origin, "http://localhost") {
		// The cookie related headers are used for the api requests authentication
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET,POST,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "cookie,content-type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		if req.Method == "OPTIONS" {
			w.WriteHeader(200)
			w.Write([]byte("OK"))
		}
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
