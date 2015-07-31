package util

import "net/http"

// CheckMethod is here to check whether an incoming request
// uses an allowed request method.
// If the request method is not one of methods, CheckMethod
// will return false.
//
// FIXME(flowlo): We have a regression here. Previously,
// CheckMethod would answer the request with an Allow header.
// This is not the case anymore.
func CheckMethod(r *http.Request, methods ...string) (ok bool) {
	for _, method := range methods {
		if method == r.Method {
			return true
		}
	}
	return
}
