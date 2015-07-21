package controllers

import (
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"github.com/coduno/app/util"
)

// StartRun starts a run for the received CodeData
func StartRun(w http.ResponseWriter, r *http.Request, c context.Context) {
	if !util.CheckMethod(w, r, "POST") {
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Error reading: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := http.Post("http://localhost:8081/api/run/start/simple", "raw", strings.NewReader(string(body)))

	if err != nil {
		http.Error(w, "Error sending code: "+err.Error(), http.StatusInternalServerError)
		return
	}

	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		http.Error(w, "Error reading: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

func DownloadTemplate(w http.ResponseWriter, r *http.Request, c context.Context) {
	if !util.CheckMethod(w, r, "GET") {
		return
	}
	// TODO(victorbalan): Serve correct template using the username and fingerprint id
	// urlParams := mux.Vars(r)
	// token := urlParams["token"]
	// username, fingerprint, err := DecodeToken(token)

	// FIXME(victorbalan): Send correct Content-Type
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename='template.java'")

	http.ServeFile(w, r, "challenges/template.java")
	return
}
