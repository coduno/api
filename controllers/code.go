package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/coduno/app/models"
	"github.com/coduno/app/util"
)

func UploadCode(w http.ResponseWriter, r *http.Request) {
	if !util.CheckMethod(w, r, "POST") {
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Error reading: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var codeData models.CodeData
	err = json.Unmarshal(body, &codeData)

	if err != nil {
		http.Error(w, "Cannot unmarshal: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// TODO(victorbalan): Pass the code to the engine

	w.Write([]byte("Success"))
}

func DownloadTemplate(w http.ResponseWriter, r *http.Request) {
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
