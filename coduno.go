package coduno

import (
	"appengine"
	"appengine/urlfetch"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var gitlabToken = "YHQiqMx3qUfj8_FxpFe4"

type ContainerCreation struct {
	Id string `json:"Id"`
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
	http.HandleFunc("/api/push", setupHandler(push))
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
