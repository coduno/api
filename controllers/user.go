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
	States []string
}

func init() {
	router.Handle("/user/company", ContextHandlerFunc(GetCompanyByUser))
	router.Handle("/user", ContextHandlerFunc(WhoAmI))
	router.Handle("/users", ContextHandlerFunc(User))
	router.Handle("/users/{key}", ContextHandlerFunc(GetUser))
	router.Handle("/users/{key}/profile", ContextHandlerFunc(GetProfileForUser))
	router.Handle("/whoami", ContextHandlerFunc(WhoAmI))
}

func User(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	switch r.Method {
	case "POST":
		return createUser(ctx, w, r)
	case "GET":
		p, ok := passenger.FromContext(ctx)
		if !ok {
			return http.StatusUnauthorized, nil
		}
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
	var users model.Users
	keys, err := model.NewQueryForUser().GetAll(ctx, &users)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	json.NewEncoder(w).Encode(users.Key(keys))
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

func GetUsersInvitedByCompany(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}
	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}
	var u model.User
	if err = datastore.Get(ctx, p.User, &u); err != nil {
		return http.StatusInternalServerError, nil
	}

	if u.Company == nil {
		return http.StatusUnauthorized, nil
	}

	var invitations model.Invitations
	_, err = model.NewQueryForInvitation().Ancestor(u.Company).GetAll(ctx, &invitations)
	m := make(map[string]int)
	for _, inv := range invitations {
		ek := inv.User.Encode()
		if _, ok := m[ek]; !ok {
			m[ek] = 1
		} else {
			m[ek] = m[ek] + 1
		}
	}

	var userList []*keyedUserWithState
	var kUser *keyedUserWithState
	for k, v := range m {
		if kUser, err = getUserState(ctx, k, v); err != nil {
			return http.StatusInternalServerError, err
		}
		userList = append(userList, kUser)
	}
	json.NewEncoder(w).Encode(userList)
	return http.StatusOK, nil
}

func getUserState(ctx context.Context, ek string, i int) (*keyedUserWithState, error) {
	key, err := datastore.DecodeKey(ek)
	if err != nil {
		return nil, err
	}
	var user model.User
	if err = datastore.Get(ctx, key, &user); err != nil {
		return nil, err
	}
	var (
		profiles model.Profiles
		keys     []*datastore.Key
	)
	if keys, err = model.NewQueryForProfile().Ancestor(key).GetAll(ctx, &profiles); err != nil {
		return nil, err
	}

	if len(keys) < 1 {
		return nil, errors.New("Profile not found")
	}

	kUser := &keyedUserWithState{
		user.Key(key),
		make([]string, 0),
	}
	var results model.Results
	if _, err := model.NewQueryForResult().
		Ancestor(keys[0]).
		Order("Started").
		GetAll(ctx, &results); err != nil {
		return nil, err
	}
	if len(results) < i {
		// User was invited but has not yet started at least one of the
		// challenges that he was invited at.
		kUser.States = append(kUser.States, "invited")
	}
	finished := true
	for _, result := range results {
		if result.Started.Before(time.Now()) && result.Finished.IsZero() {
			kUser.States = append(kUser.States, "coding")
			finished = false
			break
		}
	}
	if len(results) == i && finished {
		kUser.States = append(kUser.States, "finished")
	}

	return kUser, nil
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
