package controllers

import (
	"errors"
	"io"
	"net/http"

	"github.com/coduno/api/util"
	"github.com/gorilla/mux"

	"google.golang.org/cloud/storage"

	"golang.org/x/net/context"
)

// Template serves the contents of a static file to a client.
// TODO(flowlo, victorbalan): Decide where the templates will be stored.
func Template(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	templateName, ok := mux.Vars(r)["name"]
	if !ok {
		return http.StatusBadRequest, errors.New("template name missing")
	}

	// TODO(victorbalan): pass the template file and not just the folder
	// so we can be more flexible with templates.
	language := mux.Vars(r)["language"]
	fileName, ok := util.FileNames[language]
	if !ok {
		return http.StatusInternalServerError, errors.New("language unknown")
	}

	rdr, err := storage.NewReader(util.CloudContext(ctx), util.TemplateBucket, templateName+"/"+fileName)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename='"+fileName+"'")
	go io.Copy(w, rdr)
	return http.StatusOK, nil
}
