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
		key, err = user.Save(ctx)
	} else {
		// Bind user to company for eternity.
		key, err = user.SaveWithParent(ctx, companyKey)
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user.Key(key))
	return http.StatusOK, nil
}

func getUsers(p *passenger.Passenger, ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if p.UserKey.Parent() == nil {
		return http.StatusUnauthorized, nil
	}
	var invitations model.Invitations
	_, err = model.NewQueryForInvitation().Ancestor(p.UserKey.Parent()).GetAll(ctx, &invitations)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	cckeys, err := model.NewQueryForChallenge().Ancestor(p.UserKey.Parent()).KeysOnly().GetAll(ctx, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	var resultKeys []*datastore.Key
	for _, val := range cckeys {
		rkeys, err := model.NewQueryForResult().Filter("Challenge =", val).KeysOnly().GetAll(ctx, nil)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		resultKeys = append(resultKeys, rkeys...)
	}

	var users model.Users
	keys, err := model.NewQueryForUser().GetAll(ctx, &users)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	finishedUsers := make([]*datastore.Key, len(resultKeys))
	for i := range resultKeys {
		finishedUsers[i] = resultKeys[i].Parent().Parent()
	}

	// TODO(victorbalan): Don`t load invited users that have an result.
	invitedUsers := make([]*datastore.Key, len(invitations))
	for i, val := range invitations {
		invitedUsers[i] = val.User
	}
	mappedStates := make(map[string]string)
	for _, val := range invitedUsers {
		mappedStates[val.Encode()] = "invited"
	}
	for _, val := range finishedUsers {
		mappedStates[val.Encode()] = "coding"
	}

	usersWithState := make([]keyedUserWithState, len(users))
	for i := range users {
		usersWithState[i] = keyedUserWithState{
			KeyedUser: &model.KeyedUser{
				User: &users[i], Key: keys[i],
			},
			State: mappedStates[keys[i].Encode()],
		}
	}
	json.NewEncoder(w).Encode(usersWithState)
	return http.StatusOK, nil
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
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	var user model.User
	if err := datastore.Get(ctx, p.UserKey, &user); err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(user.Key(p.UserKey))
	return http.StatusOK, nil
}
