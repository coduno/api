package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/mail"
	"time"

	ourmail "github.com/coduno/app/mail"
	"github.com/coduno/engine/model"
	"github.com/coduno/engine/passenger"
	"github.com/coduno/engine/util"
	"github.com/coduno/engine/util/password"
	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
)

// Invitation handles the creation of a new invitation and sends an e-mail to
// the user.
func Invitation(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	p, ok := passenger.FromContext(ctx)
	if !ok {
		http.Error(w, "permission denied", http.StatusUnauthorized)
		return
	}
	cKey := p.UserKey.Parent()
	if cKey == nil {
		http.Error(w, "permission denied", http.StatusUnauthorized)
	}
	var company model.Company
	datastore.Get(ctx, cKey, &company)

	// TODO(flowlo): Also check whether the parent of the current user is the
	// parent of the challenge (if any).

	if !util.CheckMethod(w, r, "POST") {
		return
	}

	var params = struct {
		Address   string
		Challenge *datastore.Key
	}{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	address, err := mail.ParseAddress(params.Address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var users model.Users
	keys, err := model.NewQueryForUser().
		Filter("Address=", address.Address).
		KeysOnly().
		Limit(1).
		GetAll(ctx, &users)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var key *datastore.Key
	var user model.User
	if len(keys) == 1 {
		key = keys[0]
		user = users[0]
	} else {
		user = model.User{Address: *address}
		key, err = datastore.Put(ctx, key, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// TODO(flowlo): Generate token with its own util.
	tokenValue, err := password.Generate(0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	now := time.Now()
	accessToken := model.AccessToken{
		Value:        string(tokenValue),
		Creation:     now,
		Modification: now,
		Expiry:       now.Add(time.Hour * 24 * 365),
		Description:  "Initialization Token",
	}

	token := base64.URLEncoding.EncodeToString([]byte(params.Challenge.Encode() + accessToken.Value))

	invitation := model.Invitation{
		User: key,
	}

	buf := new(bytes.Buffer)
	if err := ourmail.Invitation.Execute(buf, struct {
		UserAddress, CompanyAddress mail.Address
		Token                       string
	}{
		user.Address,
		company.Address,
		token,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ourmail.Send(ctx, user.Address, "We challenge you!", buf.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key, err = invitation.Save(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	invitation.Write(w, key)
}
