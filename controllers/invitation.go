package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"time"

	ourmail "github.com/coduno/app/mail"
	"github.com/coduno/app/model"
	"github.com/coduno/app/util"
	"github.com/coduno/engine/passenger"
	"github.com/coduno/engine/util/password"
	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
)

// Invitation handles the creation of a new invitation and sends an e-mail to
// the user.
func Invitation(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, errors.New("permission denied")
	}
	cKey := p.UserKey.Parent()
	if cKey == nil {
		return http.StatusUnauthorized, errors.New("permission denied")
	}
	var company model.Company
	datastore.Get(ctx, cKey, &company)

	// TODO(flowlo): Also check whether the parent of the current user is the
	// parent of the challenge (if any).

	if err = util.CheckMethod(r, "GET"); err != nil {
		return http.StatusMethodNotAllowed, err
	}

	var params = struct {
		Address   string
		Challenge *datastore.Key
	}{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return http.StatusBadRequest, err
	}

	address, err := mail.ParseAddress(params.Address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return http.StatusBadRequest, err
	}

	var users model.Users
	keys, err := model.NewQueryForUser().
		Filter("Address=", address.Address).
		KeysOnly().
		Limit(1).
		GetAll(ctx, &users)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return http.StatusInternalServerError, err
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
			return http.StatusInternalServerError, err
		}
	}

	// TODO(flowlo): Generate token with its own util.
	tokenValue, err := password.Generate(0)
	if err != nil {
		return http.StatusInternalServerError, err
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
		return http.StatusInternalServerError, err
	}

	err = ourmail.Send(ctx, user.Address, "We challenge you!", buf.String())
	if err != nil {
		return http.StatusInternalServerError, err
	}

	key, err = invitation.Save(ctx)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	invitation.Write(w, key)
	return http.StatusOK, nil
}
