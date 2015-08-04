package controllers

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

// Certificate will return all public certificates assigned by App Engine in PEM format.
func Certificate(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	var certs []appengine.Certificate
	if certs, err = appengine.PublicCertificates(ctx); err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Add("Content-Type", "application/x-pem-file")

	for _, cert := range certs {
		w.Write([]byte("KeyName: \"" + cert.KeyName + "\""))
		w.Write(cert.Data)
	}

	return http.StatusOK, nil
}
