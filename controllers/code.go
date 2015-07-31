package controllers

import (
	"net/http"

	"golang.org/x/net/context"
)

// DownloadTemplate serves a static file to a client.
// TODO(flowlo, victorbalan): Decide where the templates will be stored.
func Template(c context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	// TODO(victorbalan): Serve correct template using passenger and result key.

	// FIXME(victorbalan): Send correct Content-Type
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename='template.java'")

	http.ServeFile(w, r, "challenges/template.java")
	return http.StatusOK, nil
}
