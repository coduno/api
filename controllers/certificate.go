package controllers

import (
	"net/http"

	"google.golang.org/appengine"
)

// Certificate will return all public certificates assigned by App Engine in PEM format.
func Certificate(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	certs, err := appengine.PublicCertificates(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/x-pem-file")

	for _, cert := range certs {
		w.Write([]byte("KeyName: \"" + cert.KeyName + "\""))
		w.Write(cert.Data)
	}
}
