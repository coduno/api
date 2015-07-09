package controllers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine/log"

	"models"
)

func Push(w http.ResponseWriter, r *http.Request, c context.Context) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Invalid method.", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Error reading: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var build models.Build
	err = json.Unmarshal(body, &build)

	if err != nil {
		http.Error(w, "Cannot unmarshal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//err = build.Create()
	//err = build.Start()

	log.Infof(c, "Received push for %s at %s.", build.Hash, build.Repository)
}

func Lol(w http.ResponseWriter, r *http.Request, c context.Context) {
	io.WriteString(w, "lol")
}
