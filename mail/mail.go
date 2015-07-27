package mail

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"text/template"

	"golang.org/x/net/context"
	appmail "google.golang.org/appengine/mail"
)

//CandidateInvitationTemplate is the name of the e-mail template for candidate invite
const CandidateInvitationTemplate = "candidate"

// SubscriptionTemplate is the name of the e-mail template for the subscription
const SubscriptionTemplate = "subscription"

var templates map[string]*template.Template

func init() {
	templates = make(map[string]*template.Template)
	templates[CandidateInvitationTemplate] = initTemplate("./mail/template.invitation")
	templates[SubscriptionTemplate] = initTemplate("./mail/template.subscription")
}

func initTemplate(path string) *template.Template {
	temp := template.New(path)
	m := make(template.FuncMap)
	m["encode"] = hex.EncodeToString
	temp.Funcs(m)

	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	message := string(dat[:])
	temp, err = temp.Parse(message)

	if err != nil {
		panic(err)
	}
	return temp
}

func SendMail(c context.Context, to []string, subject, body string) error {
	return appmail.Send(c, &appmail.Message{
		Sender:  "Lorenz Leutgeb <lorenz.leutgeb@cod.uno>",
		To:      to,
		Subject: subject,
		Body:    body,
	})
}

//PrepareMailTemplate parses the mail template defined and fills it with data
func PrepareMailTemplate(templatePath string, data interface{}) (string, error) {
	temp := template.New("template")
	m := make(template.FuncMap)
	m["encode"] = hex.EncodeToString
	temp.Funcs(m)

	// Prepare verify mail template
	dat, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", err
	}
	message := string(dat[:])
	temp, err = temp.Parse(message)
	if err != nil {
		panic(err)
	}
	var body bytes.Buffer
	err = temp.Execute(&body, data)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}
