package util

import (
	"encoding/json"
	"errors"
	"net/http"
)

// SetCookie is a simple wrapper around http.SetCookie that
// enforces HttpOnly and Secure flags.
func SetCookie(w http.ResponseWriter, cookie *http.Cookie) {
	cookie.Secure, cookie.HttpOnly = true, true
	http.SetCookie(w, cookie)
}

// CheckMethod is here to check whether an incoming request
// uses an allowed request method.
// If the request method is not one of methods, CheckMethod
// will answer the request and return false.
func CheckMethod(r *http.Request, methods ...string) error {
	for _, method := range methods {
		if method == r.Method {
			return nil
		}
	}
	return errors.New("Invalid method " + r.Method + " not allowed.")
}

// WriteMap marshals a map into json and writes it to the client
func WriteMap(w http.ResponseWriter, data map[string]interface{}) {
	body, err := json.Marshal(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(body)
}
