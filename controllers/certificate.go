package controllers

import (
	"net/http"

	"google.golang.org/appengine"
)

func init() {
	router.Handle("/cert", hsts(Certificate))
}

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
		w.Write([]byte("KeyName: \"" + cert.KeyName + "\"\n"))
		w.Write(cert.Data)
		w.Write([]byte("\n"))
	}
}
