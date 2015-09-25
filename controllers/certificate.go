package controllers

import (
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
)

func init() {
	router.Handle("/cert", ContextHandlerFunc(Certificate))
}

// Certificate will return all public certificates assigned by App Engine in PEM format.
func Certificate(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	certs, err := appengine.PublicCertificates(ctx)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// TODO(flowlo): Is this really the appropriate MIME type?
	w.Header().Add("Content-Type", "application/x-pem-file")

	for _, cert := range certs {
		p := make([]byte, 13+len(cert.KeyName)+len(cert.Data))
		p = append(p, []byte("KeyName: \""+cert.KeyName+"\"\n")...)
		p = append(p, cert.Data...)
		p = append(p, []byte("\n")...)
		if _, err := w.Write(p); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}
