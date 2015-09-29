package controllers

import (
	"errors"
	"io"
	"net/http"

	"github.com/coduno/api/util"
	"github.com/gorilla/mux"

	"golang.org/x/net/context"
)

func init() {
	router.Handle("/code/download/{name}/{language}", ContextHandlerFunc(Template))
}

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
	fn := templateName + "/" + fileName
	rc, err := util.Load(ctx, util.TemplateBucket, fn)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer rc.Close()
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename='"+fileName+"'")
	io.Copy(w, rc)
	return http.StatusOK, nil
}
