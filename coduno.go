package coduno

import (
	"appengine"
	"appengine/datastore"
	appmail "appengine/mail"
	"appengine/urlfetch"
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gitlab"
	"io/ioutil"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"time"
)

var gitlabToken = "YHQiqMx3qUfj8_FxpFe4"

type ContainerCreation struct {
	Id string `json:"Id"`
}

func connectDatabase() (*sql.DB, error) {
	return sql.Open("mysql", "root@cloudsql(coduno:mysql)/coduno")
}

type Handler func(http.ResponseWriter, *http.Request)

func setupHandler(handler Handler) Handler {
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

		handler(w, r)
	}
}

func init() {
	http.HandleFunc("/api/token", setupHandler(token))
	http.HandleFunc("/push", setupHandler(push))
	http.HandleFunc("/subscriptions", setupHandler(subscriptions))
}

type Subscription struct {
	Address          string
	EntryTime        time.Time
	Token            []byte
	VerificationTime time.Time
}

func subscriptions(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		subscribe(w, r)
	} else if r.Method == "GET" {
		action := r.FormValue("action")
		if action == "confirm" {
			confirm(w, r)
		} else if action == "delete" {
			delete(w, r)
		} else {
			http.Error(w, "Unknown action.", http.StatusBadRequest)
		}
	} else {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	address, err := mail.ParseAddress(r.FormValue("email"))

	if err != nil {
		http.Error(w, "Invalid email address: "+err.Error(), 422)
		return
	}

	token, err := hex.DecodeString(r.FormValue("token"))

	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), 422)
		return
	}

	c := appengine.NewContext(r)
	q := datastore.NewQuery("subscription").Filter("Address =", address.Address).Limit(1)

	var subs []Subscription
	keys, err := q.GetAll(c, &subs)

	if len(subs) != 1 {
		http.Error(w, "Unable to identify your subscription. Got "+fmt.Sprintf("%d", len(subs))+" matches on "+address.Address, http.StatusInternalServerError)
		return
	}

	sub := subs[0]

	if !bytes.Equal(sub.Token, token) {
		http.Error(w, "Forbidden", http.StatusUnauthorized)
		return
	}

	err = datastore.Delete(c, keys[0])

	if err != nil {
		http.Error(w, "Failed to delete your subscription. Please contact root@cod.uno.", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Your subscription was removed, we're sorry to see you go :("))
}

func confirm(w http.ResponseWriter, r *http.Request) {
	address, err := mail.ParseAddress(r.FormValue("email"))

	if err != nil {
		http.Error(w, "Invalid email address: "+err.Error(), 422)
		return
	}

	token, err := hex.DecodeString(r.FormValue("token"))

	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), 422)
		return
	}

	c := appengine.NewContext(r)
	q := datastore.NewQuery("subscription").Filter("Address =", address.Address).Limit(1)

	var subs []Subscription
	keys, err := q.GetAll(c, &subs)

	if len(subs) != 1 {
		http.Error(w, "Unable to identify your subscription. Got "+fmt.Sprintf("%d", len(subs))+" matches on "+address.Address, http.StatusInternalServerError)
		return
	}

	sub := subs[0]

	if !bytes.Equal(sub.Token, token) {
		http.Error(w, "Forbidden", http.StatusUnauthorized)
		return
	}

	sub.VerificationTime = time.Now()
	token, _, err = newToken()

	if err != nil {
		http.Error(w, "Unable to generate your new token.", http.StatusInternalServerError)
		return
	}

	sub.Token = token

	_, err = datastore.Put(c, keys[0], &sub)

	if err != nil {
		http.Error(w, "Failed to store subscription: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Your subscription was verified, yay!"))
}

func newToken() ([]byte, int, error) {
	token := make([]byte, 16)
	n, err := rand.Read(token)
	return token, n, err
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	address, err := mail.ParseAddress(r.FormValue("email"))

	if err != nil {
		http.Error(w, "Invalid email address: "+err.Error(), 422)
		return
	}

	q := datastore.NewQuery("subscription").Filter("Address =", address.Address).Limit(1).KeysOnly()
	collisions, err := q.GetAll(c, nil)

	if err != nil {
		http.Error(w, "Failed to check for duplicates.", http.StatusInternalServerError)
		return
	}

	if len(collisions) > 0 {
		http.Error(w, "Duplicate email address.", 422)
		return
	}

	revocationBytes, _, err := newToken()

	if err != nil {
		http.Error(w, "Failed to generate revocation secret: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sub := Subscription{
		Address:   address.Address,
		EntryTime: time.Now(),
		Token:     revocationBytes,
	}

	_, err = datastore.Put(c, datastore.NewIncompleteKey(c, "subscription", nil), &sub)

	if err != nil {
		http.Error(w, "Failed to store subscription: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = sub.RequestConfirmation(c)

	if err != nil {
		http.Error(w, "Failed to send confirmation request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("A message to confirm your subscription was sent."))
}

func (sub Subscription) RequestConfirmation(c appengine.Context) error {
	return appmail.Send(c, &appmail.Message{
		Sender:  "Lorenz Leutgeb <lorenz.leutgeb@cod.uno>",
		To:      []string{sub.Address},
		Subject: "Hello from Coduno",
		Body:    "Hey there,\r\n\r\nclick this Link to make it happen: https://www.cod.uno/subscriptions?action=confirm&email="+url.QueryEscape(sub.Address)+"&token="+hex.EncodeToString(sub.Token)+"\r\n\r\nIn case you want to undo all previous actions, please click here: https://www.cod.uno/subscriptions?action=delete&email="+url.QueryEscape(sub.Address)+"&token="+hex.EncodeToString(sub.Token),
	})
}

func push(w http.ResponseWriter, req *http.Request) {
	context := appengine.NewContext(req)
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		context.Warningf(err.Error())
	} else if len(body) < 1 {
		context.Warningf("Received empty body.")
	} else {
		push, err := gitlab.NewPush(body)

		if err != nil {
			context.Warningf(err.Error())
		} else {
			commit := push.Commits[0]

			docker, _ := http.NewRequest("POST", "http://docker.cod.uno:2375/v1.15/containers/create", strings.NewReader(`
				{
					"Image": "coduno/git:experimental",
					"Cmd": ["/start.sh"],
					"Env": [
						"CODUNO_REPOSITORY_NAME=`+push.Repository.Name+`",
						"CODUNO_REPOSITORY_URL=`+push.Repository.URL+`",
						"CODUNO_REPOSITORY_HOMEPAGE=`+push.Repository.Homepage+`",
						"CODUNO_REF=`+push.Ref+`"
					]
				}
			`))

			docker.Header = map[string][]string{
				"Content-Type": {"application/json"},
				"Accept":       {"application/json"},
			}

			client := urlfetch.Client(context)
			res, err := client.Do(docker)

			if err != nil {
				context.Debugf("Docker API response:", err.Error())
				return
			}

			var result ContainerCreation
			body, err := ioutil.ReadAll(res.Body)
			err = json.Unmarshal(body, &result)

			if err != nil {
				context.Debugf("Received body %d: %s", res.StatusCode, string(body))
				context.Debugf("Unmarshalling API response: %s", err.Error())
				return
			}

			docker, _ = http.NewRequest("POST", "http://docker.cod.uno:2375/v1.15/containers/"+result.Id+"/start", nil)

			docker.Header = map[string][]string{
				"Content-Type": {"application/json"},
				"Accept":       {"application/json"},
			}

			res, err = client.Do(docker)

			if err != nil {
				context.Debugf("Docker API response 2: %s", err.Error())
			} else {
				result, err := ioutil.ReadAll(res.Body)
				if err != nil {
					context.Debugf(err.Error())
				} else {
					context.Debugf(string(result))
				}
			}
			defer res.Body.Close()

			context.Infof("Received push from %s", commit.Author.Email)
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

func token(w http.ResponseWriter, req *http.Request) {
	context := appengine.NewContext(req)

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

			client := urlfetch.Client(context)
			res, err := client.Do(gitlabReq)
			if err != nil {
				context.Debugf(err.Error())
			} else {
				result, err := ioutil.ReadAll(res.Body)
				if err != nil {
					context.Debugf(err.Error())
				} else {
					context.Debugf(string(result))
				}
			}
			defer res.Body.Close()

			context.Infof("Generated token for '%s'", username)
			fmt.Fprintf(w, token)
		}
	} else {
		// This could be either invalid/missing credentials or
		// the database failing, so let's issue a warning.
		context.Warningf(err.Error())
		http.Error(w, "Invalid credentials.", http.StatusForbidden)
	}
}

func authenticate(req *http.Request) (error, string) {
	username, password, ok := basicAuth(req)
	if !ok {
		return errors.New("No Authorization header present"), ""
	}
	return check(username, password), username
}

func check(username, password string) error {
	db, err := connectDatabase()

	if err != nil {
		return err
	}

	defer db.Close()

	rows, err := db.Query("select count(*) from users where username = ? and password = sha2(concat(?, salt), 512)", username, password)

	if err != nil {
		return err
	}

	defer rows.Close()

	var result string
	rows.Next()
	rows.Scan(&result)

	if result != "1" {
		return errors.New("Failed to validate credentials.")
	}

	return nil
}

// BasicAuth returns the username and password provided in the request's
// Authorization header, if the request uses HTTP Basic Authentication.
// See RFC 2617, Section 2.
// This will be obsoleted by go1.4
func basicAuth(r *http.Request) (username, password string, ok bool) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return
	}
	return parseBasicAuth(auth)
}

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
// This will be obsoleted by go1.4
func parseBasicAuth(auth string) (username, password string, ok bool) {
	if !strings.HasPrefix(auth, "Basic ") {
		return
	}
	c, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}
