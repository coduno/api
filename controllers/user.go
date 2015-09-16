package controllers

import (
	"errors"
	"net/http"
	"net/mail"
	"time"

	"encoding/json"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"github.com/coduno/api/util/passenger"
	"github.com/coduno/api/util/password"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type keyedUserWithState struct {
	*model.KeyedUser
	State string
}

func User(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}
	switch r.Method {
	case "POST":
		return createUser(ctx, w, r)
	case "GET":
		return getUsers(p, ctx, w, r)
	default:
		return http.StatusMethodNotAllowed, nil
	}
}

func GetUser(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	_, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	var userKey *datastore.Key
	if userKey, err = datastore.DecodeKey(mux.Vars(r)["key"]); err != nil {
		return http.StatusInternalServerError, err
	}

	var user model.User
	if err = datastore.Get(ctx, userKey, &user); err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(user.Key(userKey))
	return
}

func createUser(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	var body = struct {
		Address, Nick, Password, Company string
	}{}

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return http.StatusBadRequest, err
	}

	var companyKey *datastore.Key
	if body.Company != "" {
		companyKey, err = datastore.DecodeKey(body.Company)
		if err != nil {
			return http.StatusBadRequest, err
		}
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

	var key *datastore.Key
	if companyKey == nil {
		key, err = user.Put(ctx, nil)
	} else {
		// Bind user to company for eternity.
		key, err = user.PutWithParent(ctx, companyKey)
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user.Key(key))
	return http.StatusOK, nil
}

func getUsers(p *passenger.Passenger, ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	var u model.User
	if err = datastore.Get(ctx, p.User, &u); err != nil {
		return http.StatusInternalServerError, nil
	}

	if u.Company == nil {
		return http.StatusUnauthorized, nil
	}
	var invitations model.Invitations
	_, err = model.NewQueryForInvitation().Ancestor(u.Company).GetAll(ctx, &invitations)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	var state string
	invitedUsers := make(map[*datastore.Key]string)
	for _, invitation := range invitations {
		if state, err = getUserState(ctx, invitation); err != nil {
			return http.StatusInternalServerError, err
		}
		invitedUsers[invitation.User] = state
	}

	var usersWithState []keyedUserWithState
	var user model.User
	for key, state := range invitedUsers {

		if err = datastore.Get(ctx, key, &user); err != nil {
			return http.StatusInternalServerError, err
		}
		usersWithState = append(usersWithState, keyedUserWithState{
			KeyedUser: &model.KeyedUser{
				User: &user,
				Key:  key,
			},
			State: state,
		})
	}
	json.NewEncoder(w).Encode(usersWithState)
	return http.StatusOK, nil
}

func getUserState(ctx context.Context, invitation model.Invitation) (string, error) {
	var profiles model.Profiles
	var results model.Results
	keys, err := model.NewQueryForProfile().Ancestor(invitation.User).GetAll(ctx, &profiles)
	// Profile should have been created when the user was invited to perform the challenge.
	if err != nil || len(keys) != 1 {
		return "", err
	}
	if _, err := model.NewQueryForResult().Ancestor(keys[0]).GetAll(ctx, &results); err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "invited", nil
	}
	if results[0].Started.Before(time.Now()) && results[0].Finished.IsZero() {
		return "coding", nil
	}
	return "finished", nil
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
	if r.Method != "GET" {
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
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	var u model.User
	if err := datastore.Get(ctx, p.User, &u); err != nil {
		return http.StatusInternalServerError, nil
	}

	if u.Company == nil {
		return http.StatusUnauthorized, nil
	}

	// The account is associated with a company, so we return it.
	var company model.Company
	if err := datastore.Get(ctx, u.Company, &company); err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(company.Key(u.Company))
	return http.StatusOK, nil
}

func WhoAmI(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	var user model.User
	if err := datastore.Get(ctx, p.User, &user); err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(user.Key(p.User))
	return http.StatusOK, nil
}
