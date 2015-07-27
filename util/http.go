package util

import (
	"encoding/json"
	"net/http"
	"strings"

	"google.golang.org/appengine/datastore"
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
func CheckMethod(w http.ResponseWriter, r *http.Request, methods ...string) (ok bool) {
	for _, method := range methods {
		if method == r.Method {
			return true
		}
	}
	w.Header().Set("Allow", strings.Join(methods, ", "))
	http.Error(w, "Invalid method "+r.Method+" not allowed.", http.StatusMethodNotAllowed)
	return
}

// WriteEntities takes keys and values, generates a map
// from them and writes JSON to w.
// len(keys) != len(values) or an error during marshaling
// will result in an internal server error.
func WriteEntities(w http.ResponseWriter, keys []*datastore.Key, values []interface{}) {
	if len(keys) != len(values) {
		http.Error(w, "length mismatch while writing entities", http.StatusInternalServerError)
		return
	}

	tmp := make(map[string]interface{}, len(keys))
	for i := 0; i < len(keys); i++ {
		tmp[keys[i].String()] = values[i]
	}

	WriteMap(w, tmp)
}

// WriteEntity takes a key and the corresponding entity and writes
// it out to w after marshaling to JSON.
func WriteEntity(w http.ResponseWriter, key *datastore.Key, value interface{}) {
	WriteMap(w, map[string]interface{}{
		key.Encode(): value,
	})
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
