package controllers

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/mail"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
	appmail "google.golang.org/appengine/mail"
)

var echoMailFunc = delay.Func("echoMail", echoMail)

func echoMail(ctx context.Context, m mail.Message) {
	from, err := m.Header.AddressList("From")
	if err != nil {
		log.Warningf(ctx, "Failed getting sender of mail: %+v", m)
		return
	}

	b, _ := ioutil.ReadAll(m.Body)

	am := &appmail.Message{
		Sender:  "lorenz.leutgeb@cod.uno",
		To:      []string{from[0].String()},
		Body:    string(b),
		Headers: m.Header,
	}

	err = appmail.Send(ctx, am)
	if err != nil {
		log.Errorf(ctx, err.Error())
	}
}

// ReceiveMail will receive an e-mail and echo it back to the sender.
func ReceiveMail(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	m, err := mail.ReadMessage(r.Body)
	if err != nil {
		log.Errorf(ctx, "Failed reading a mail!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = echoMailFunc.Call(ctx, m)
	if err != nil {
		log.Errorf(ctx, "Failed enqueing handler for a mail!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, "OK")

	// TODO(flowlo):
	//  1. Check whether range m.Header.AddressList("From")
	//     fits a registered customer
	//  2. Filter mail content for further e-mail addresses
	//  3. Create a Fingerprint
	//  4. Mail the Fingerprint URL to the other address
}
