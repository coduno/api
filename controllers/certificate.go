package controllers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

func Certificate(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	var c []appengine.Certificate
	if c, err = appengine.PublicCertificates(ctx); err != nil {
		return http.StatusInternalServerError, err
	}

	if err = json.NewEncoder(w).Encode(c); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
