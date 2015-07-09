package util

import (
	"net/http"
	"strings"
)

// SetCookie is a simple wrapper around http.SetCookie that
// enforces HttpOnly and Secure flags.
func SetCookie(w http.ResponseWriter, cookie *http.Cookie) {
	cookie.Secure, cookie.HttpOnly = true, true
	http.SetCookie(w, cookie)
}

func CheckMethod(w http.ResponseWriter, r *http.Request, methods ...string) bool {
	for _, method := range methods {
		if method == r.Method {
			return true
		}
	}
	w.Header().Set("Allow", strings.Join(methods, ", "))
	http.Error(w, "Invalid method "+r.Method+" not allowed.", http.StatusMethodNotAllowed)
	return false
}
