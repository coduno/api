package main

import (
	"fmt"
	"io"
	"net/http"
	"net/mail"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func receiveMail(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	m, err := mail.ReadMessage(r.Body)
	if err != nil {
		log.Warningf(ctx, "Failed reading an mail!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	from, err := m.Header.AddressList("From")
	if err != nil {
		log.Warningf(ctx, "Failed getting sender of mail: %+v", m)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("Received mail from %+v", from)
	log.Debugf(appengine.NewContext(r), response)
	io.WriteString(w, response)

	// TODO(flowlo):
	//  1. Check whether range m.Header.AddressList("From")
	//     fits a registered customer
	//  2. Filter mail content for further e-mail addresses
	//  3. Create a Fingerprint
	//  4. Mail the Fingerprint URL to the other address
}
