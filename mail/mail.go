package mail

import (
	"appengine"
	"appengine/datastore"
	appmail "appengine/mail"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"
	"text/template"
	"time"
)

type Subscription struct {
	Address          string
	EntryTime        time.Time
	Token            []byte
	VerificationTime time.Time
}

var verify = template.New("verify")

func init() {
	// Add function map to template engine
	m := make(template.FuncMap)
	m["encode"] = hex.EncodeToString
	verify.Funcs(m)

	// Prepare verify mail template
	dat, err := ioutil.ReadFile("./mail/mail.verify")
	if err != nil {
		panic(err)
	}
	message := string(dat[:])
	verify, err = verify.Parse(message)
	if err != nil {
		panic(err)
	}
}

func Subscriptions(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		subscribe(w, r)
	} else if r.Method == "GET" {
		action := r.FormValue("action")
		if action == "confirm" {
			confirm(w, r)
		} else if action == "delete" {
			delete(w, r)
		} else {
			c := appengine.NewContext(r)
			c.Infof("What happend?")
			http.Error(w, "Unknown action.", http.StatusBadRequest)
		}
	} else {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	address, err := mail.ParseAddress(r.FormValue("email"))

	if err != nil {
		http.Error(w, "Invalid email address: "+err.Error(), 422)
		return
	}

	token, err := hex.DecodeString(r.FormValue("token"))

	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), 422)
		return
	}

	c := appengine.NewContext(r)
	q := datastore.NewQuery("subscription").Filter("Address =", address.Address).Limit(1)

	var subs []Subscription
	keys, err := q.GetAll(c, &subs)

	if len(subs) != 1 {
		http.Error(w, "Unable to identify your subscription. Got "+fmt.Sprintf("%d", len(subs))+" matches on "+address.Address, http.StatusInternalServerError)
		return
	}

	sub := subs[0]

	if !bytes.Equal(sub.Token, token) {
		http.Error(w, "Forbidden", http.StatusUnauthorized)
		return
	}

	err = datastore.Delete(c, keys[0])

	if err != nil {
		http.Error(w, "Failed to delete your subscription. Please contact root@cod.uno.", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Your subscription was removed, we're sorry to see you go :("))
}

func confirm(w http.ResponseWriter, r *http.Request) {
	address, err := mail.ParseAddress(r.FormValue("email"))

	if err != nil {
		http.Error(w, "Invalid email address: "+err.Error(), 422)
		return
	}

	token, err := hex.DecodeString(r.FormValue("token"))

	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), 422)
		return
	}

	c := appengine.NewContext(r)
	q := datastore.NewQuery("subscription").Filter("Address =", address.Address).Limit(1)

	var subs []Subscription
	keys, err := q.GetAll(c, &subs)

	if len(subs) != 1 {
		http.Error(w, "Unable to identify your subscription. Got "+fmt.Sprintf("%d", len(subs))+" matches on "+address.Address, http.StatusInternalServerError)
		return
	}

	sub := subs[0]

	if !bytes.Equal(sub.Token, token) {
		http.Error(w, "Forbidden", http.StatusUnauthorized)
		return
	}

	sub.VerificationTime = time.Now()
	token, _, err = newToken()

	if err != nil {
		http.Error(w, "Unable to generate your new token.", http.StatusInternalServerError)
		return
	}

	sub.Token = token

	_, err = datastore.Put(c, keys[0], &sub)

	if err != nil {
		http.Error(w, "Failed to store subscription: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Your subscription was verified, yay!"))
}

func newToken() ([]byte, int, error) {
	token := make([]byte, 16)
	n, err := rand.Read(token)
	return token, n, err
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	address, err := mail.ParseAddress(r.FormValue("email"))

	if err != nil {
		http.Error(w, "Invalid email address: "+err.Error(), 422)
		return
	}

	q := datastore.NewQuery("subscription").Filter("Address =", address.Address).Limit(1).KeysOnly()
	collisions, err := q.GetAll(c, nil)

	if err != nil {
		http.Error(w, "Failed to check for duplicates.", http.StatusInternalServerError)
		return
	}

	if len(collisions) > 0 {
		http.Error(w, "Duplicate email address.", 422)
		return
	}

	revocationBytes, _, err := newToken()

	if err != nil {
		http.Error(w, "Failed to generate revocation secret: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sub := Subscription{
		Address:   address.Address,
		EntryTime: time.Now(),
		Token:     revocationBytes,
	}

	_, err = datastore.Put(c, datastore.NewIncompleteKey(c, "subscription", nil), &sub)

	if err != nil {
		http.Error(w, "Failed to store subscription: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = sub.RequestConfirmation(c)

	if err != nil {
		http.Error(w, "Failed to send confirmation request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("A message to confirm your subscription was sent."))
}

func (sub Subscription) RequestConfirmation(c appengine.Context) error {
	var message bytes.Buffer
	err := verify.Execute(&message, sub)
	if err != nil {
		return err
	}

	return appmail.Send(c, &appmail.Message{
		Sender:  "Lorenz Leutgeb <lorenz.leutgeb@cod.uno>",
		To:      []string{sub.Address},
		Subject: "Hello from Coduno",
		Body:    message.String(),
	})
}
