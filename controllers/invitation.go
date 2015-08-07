package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/mail"
	"text/template"
	"time"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util/passenger"
	"google.golang.org/appengine/datastore"
	appmail "google.golang.org/appengine/mail"

	"golang.org/x/net/context"
)

var invitation *template.Template

func init() {
	var err error
	invitation, err = template.ParseFiles("./mail/template.invitation")
	if err != nil {
		panic(err)
	}
}

// Invitation handles the creation of a new invitation and sends an e-mail to
// the user.
func Invitation(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, nil
	}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	cKey := p.UserKey.Parent()
	if cKey == nil {
		return http.StatusUnauthorized, nil
	}

	var company model.Company
	if err = datastore.Get(ctx, cKey, &company); err != nil {
		return http.StatusInternalServerError, err
	}

	var params = struct {
		Address, Challenge string
	}{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return http.StatusBadRequest, err
	}

	address, err := mail.ParseAddress(params.Address)
	if err != nil {
		return http.StatusBadRequest, err
	}

	challengeKey, err := datastore.DecodeKey(params.Challenge)
	if err != nil {
		return http.StatusBadRequest, err
	}

	var challenge model.Challenge
	if err := datastore.Get(ctx, challengeKey, &challenge); err != nil {
		// TODO(flowlo): Actually look into err. If it is just something like
		// "not found", an internal server error is not appropriate.
		return http.StatusInternalServerError, err
	}

	// TODO(flowlo): Check whether the parent of the current user is the
	// parent of the challenge (if any), and check whether the challenge
	// even exists.

	var users model.Users
	keys, err := model.NewQueryForUser().
		Filter("Address=", address.Address).
		KeysOnly().
		Limit(1).
		GetAll(ctx, &users)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	var key *datastore.Key
	var user model.User
	if len(keys) == 1 {
		key = keys[0]
		user = users[0]
	} else {
		user = model.User{Address: *address}
		key, err = user.Save(ctx)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		profile := model.Profile{}
		if _, err = profile.SaveWithParent(ctx, key); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	// NOTE: We are creating a new, orphaned Passenger here, because a
	// Passenger can only issue tokens for the encapsulated user.
	np := passenger.Passenger{
		UserKey: key,
	}

	now := time.Now()
	token := &model.AccessToken{
		Creation:    now,
		Expiry:      now.Add(time.Hour * 24 * 365),
		Description: "Initialization Token",
	}

	value, err := np.IssueToken(ctx, token)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	query := base64.URLEncoding.EncodeToString([]byte(params.Challenge + ":" + value))

	i := model.Invitation{
		User: key,
	}

	buf := new(bytes.Buffer)
	if err = invitation.Execute(buf, struct {
		UserAddress, CompanyAddress mail.Address
		Token                       string
	}{
		user.Address,
		company.Address,
		query,
	}); err != nil {
		return http.StatusInternalServerError, err
	}

	if err = appmail.Send(ctx, &appmail.Message{
		Sender:  "Lorenz Leutgeb <lorenz.leutgeb@cod.uno>",
		To:      []string{user.Address.String()},
		Subject: "We challenge you!",
		Body:    buf.String(),
	}); err != nil {
		return http.StatusInternalServerError, err
	}

	key, err = i.SaveWithParent(ctx, cKey)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(i.Key(key))
	return http.StatusOK, nil
}
