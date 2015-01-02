package coduno

import (
	"appengine"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strings"
)

var db = new(sql.DB)

func init() {
	db, _ = sql.Open("mysql", "root@cloudsql(coduno:db)/coduno")

	http.HandleFunc("/api/token", token)
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

	if ok, username, _ := authenticate(req); ok {
		token, err := generateToken()
		if err != nil {
			http.Error(w, "Failed to generate token.", http.StatusInternalServerError)
		} else {
			context.Infof("Generated token for '%s'", username)
			fmt.Fprintf(w, token)
		}
	} else {
		http.Error(w, "Invalid credentials.", http.StatusForbidden)
	}
}

func authenticate(req *http.Request) (bool, string, string) {
	username, password, ok := basicAuth(req)
	if !ok { // no Authorization header present
		return false, "", ""
	}
	return check(username, password)
}

func check(username, password string) (bool, string, string) {
	// TODO use DB to check for something meaningful
	return username == "flowlo" && password == "cafebabe", username, password
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
