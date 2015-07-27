package controllers

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/coduno/engine/util"
)

func DownloadTemplate(c context.Context, w http.ResponseWriter, r *http.Request) {
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
}
