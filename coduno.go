package coduno

import (
	"appengine"
	"appengine/urlfetch"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gitlab"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var gitlabToken = "YHQiqMx3qUfj8_FxpFe4"

type ContainerCreation struct {
	Id string `json:"Id"`
}

func connectDatabase() (*sql.DB, error) {
	return sql.Open("mysql", "root@cloudsql(coduno:mysql)/coduno")
}

func init() {
	http.HandleFunc("/api/token", token)
	http.HandleFunc("/push", push)
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

			docker, _ := http.NewRequest("POST", "http://docker.cod.uno:2375/v1.15/containers/create", strings.NewReader("{\"Image\": \"ubuntu\", \"Cmd\": [\"date\"]}"))

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
