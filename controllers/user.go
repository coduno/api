package controllers

import (
	"io"
	"net/http"

	"google.golang.org/appengine"
)

func User(w http.ResponseWriter, req *http.Request) {
	if appengine.IsDevAppServer() {
		io.WriteString(w, "http://localhost:8081")
	} else {
		io.WriteString(w, "https://engine.cod.uno")
	}
}
