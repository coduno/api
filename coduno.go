package coduno

import (
	"appengine"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strings"
)

func connectDatabase() (*sql.DB, error) {
	return sql.Open("mysql", "root@cloudsql(coduno:mysql)/coduno")
}

func init() {
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

	if err, username := authenticate(req); err == nil {
		token, err := generateToken()
		if err != nil {
			http.Error(w, "Failed to generate token.", http.StatusInternalServerError)
		} else {
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
