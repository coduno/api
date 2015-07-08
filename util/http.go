package util

import "net/http"

// SetCookie is a simple wrapper around http.SetCookie that
// enforces HttpOnly and Secure flags.
func SetCookie(w http.ResponseWriter, cookie *http.Cookie) {
	cookie.Secure, cookie.HttpOnly = true, true
	http.SetCookie(w, cookie)
}
