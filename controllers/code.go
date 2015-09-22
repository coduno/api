package controllers

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/coduno/api/util"
	"github.com/gorilla/mux"

	"google.golang.org/appengine/memcache"
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
	fn := templateName + "/" + fileName
	var content []byte
	var err error
	if content, err = fromCache(ctx, fn); err != nil {
		if content, err = fromStorage(ctx, fn); err != nil {
			return http.StatusInternalServerError, err
		}
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename='"+fileName+"'")
	w.Write(content)
	return http.StatusOK, nil
}

func fromCache(ctx context.Context, fileName string) ([]byte, error) {
	item, err := memcache.Get(ctx, fileName)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

func fromStorage(ctx context.Context, fileName string) ([]byte, error) {
	rc, err := storage.NewReader(util.CloudContext(ctx), util.TemplateBucket, fileName)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	item := &memcache.Item{
		Key:   fileName,
		Value: content,
	}

	if err = memcache.Set(ctx, item); err != nil {
		return nil, err
	}

	return content, nil
}
