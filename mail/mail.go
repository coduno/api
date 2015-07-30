package mail

import (
	"net/mail"
	"text/template"

	"golang.org/x/net/context"
	appmail "google.golang.org/appengine/mail"
)

var Subscription, Invitation *template.Template

func init() {
	var err error
	Subscription, err = template.ParseFiles("./mail/template.invitation")
	if err != nil {
		panic(err)
	}
	Invitation, err = template.ParseFiles("./mail/template.subscription")
	if err != nil {
		panic(err)
	}
}

func Send(c context.Context, to mail.Address, subject, body string) error {
	return appmail.Send(c, &appmail.Message{
		Sender:  "Lorenz Leutgeb <lorenz.leutgeb@cod.uno>",
		To:      []string{to.String()},
		Subject: subject,
		Body:    body,
	})
}
