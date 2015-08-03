package controllers

import (
	"errors"
	"net/http"
	"net/mail"

	"encoding/json"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"github.com/coduno/api/util/passenger"
	"github.com/coduno/api/util/password"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func User(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, nil
	}

	var body = struct {
		Address, Nick, Password string
	}{}

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return http.StatusBadRequest, err
	}

	if err = util.CheckNick(body.Nick); err != nil {
		return http.StatusBadRequest, err
	}

	var address *mail.Address
	if address, err = mail.ParseAddress(body.Address); err != nil {
		return http.StatusBadRequest, err
	}

	// Duplicate length check. If we move this after the conflict checks,
	// we could end up returning with a short password after querying Datastore.
	// The other way round, we would have to hash the password, and then throw it
	// away because of possible conflicts.
	pw := []byte(body.Password)
	if err = password.CheckLen(pw); err != nil {
		return http.StatusBadRequest, err
	}

	var emailConflict bool
	if emailConflict, err = alreadyExists(ctx, "Address", address.Address); err != nil {
		return http.StatusInternalServerError, err
	}
	if emailConflict {
		return http.StatusConflict, errors.New("duplicate e-mail address")
	}

	var nickConflict bool
	if nickConflict, err = alreadyExists(ctx, "Nick", body.Nick); err != nil {
		return http.StatusInternalServerError, err
	}
	if nickConflict {
		return http.StatusConflict, errors.New("duplicate nick")
	}

	var hashedPassword []byte
	if hashedPassword, err = password.Hash(pw); err != nil {
		return http.StatusInternalServerError, err
	}

	user := model.User{
		Address:        *address,
		Nick:           body.Nick,
		HashedPassword: hashedPassword,
	}

	key, err := user.Save(ctx)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user.Key(key))
	return
}

func alreadyExists(ctx context.Context, property, value string) (exists bool, err error) {
	k, err := model.NewQueryForUser().
		KeysOnly().
		Limit(1).
		Filter(property+"=", value).
		GetAll(ctx, nil)

	if err != nil {
		return false, err
	}

	return len(k) == 1, nil
}

// GetUsersByCompany queries the user accounts belonging to a company.
func GetUsersByCompany(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if !util.CheckMethod(r, "GET") {
		return http.StatusMethodNotAllowed, nil
	}
	key, err := datastore.DecodeKey(r.URL.Query()["result"][0])
	if err != nil {
		return http.StatusBadRequest, err
	}
	var users model.Users
	keys, err := model.NewQueryForUser().
		Ancestor(key).
		GetAll(ctx, &users)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	json.NewEncoder(w).Encode(users.Key(keys))
	return http.StatusOK, nil
}

func GetCompanyByUser(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if !util.CheckMethod(r, "GET") {
		return http.StatusMethodNotAllowed, nil
	}
	p, ok := passenger.FromContext(ctx)

	if !ok {
		return http.StatusUnauthorized, nil
	}
	key := p.UserKey.Parent()
	if key == nil {
		return http.StatusUnauthorized, nil
	}
	// The account is associated with a company, so we return it.
	var company model.Company
	if err := datastore.Get(ctx, key, &company); err != nil {
		return http.StatusInternalServerError, err
	}
	json.NewEncoder(w).Encode(company.Key(key))
	return http.StatusOK, nil
}

func WhoAmI(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if !util.CheckMethod(r, "GET") {
		return http.StatusMethodNotAllowed, nil
	}
	p, ok := passenger.FromContext(ctx)

	if !ok {
		return http.StatusUnauthorized, nil
	}
	json.NewEncoder(w).Encode(p.User.Key(p.UserKey))
	return http.StatusOK, nil
}
