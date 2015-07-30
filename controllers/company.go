package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/app/model"
	"github.com/coduno/app/util"
	"github.com/coduno/engine/util/password"
)

// LoginInfo is the login info for a company
type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CompanyLogin starts a session for a company
func CompanyLogin(c context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if err = util.CheckMethod(r, "GET"); err != nil {
		return http.StatusMethodNotAllowed, err
	}
	var loginInfo LoginInfo
	if err = json.NewDecoder(r.Body).Decode(&loginInfo); err != nil {
		return http.StatusBadRequest, err
	}

	q := model.NewQueryForUser().
		Filter("Address = ", loginInfo.Email).
		Limit(1)

	var users model.Users
	keys, err := q.GetAll(c, &users)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if len(users) != 1 {
		// NOTE: Do not leak len(users) here.
		return http.StatusUnauthorized, errors.New("Unauthorized")
	}

	user := users[0]
	key := keys[0]

	if err = password.Check([]byte(loginInfo.Password), user.HashedPassword); err != nil {
		// NOTE: Do not leak err here.
		return http.StatusUnauthorized, errors.New("Unauthorized")
	}

	user.Write(w, key)
	return http.StatusOK, nil
}

func PostCompany(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if err = util.CheckMethod(r, "GET"); err != nil {
		return http.StatusMethodNotAllowed, err
	}

	var company model.Company
	if err = json.NewDecoder(r.Body).Decode(&company); err != nil {
		return http.StatusBadRequest, err
	}

	q := model.NewQueryForCompany().
		Filter("Address = ", company.Address.Address).
		Limit(1)

	var companies model.Companys
	if _, err = q.GetAll(ctx, &companies); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(companies) > 0 {
		return http.StatusConflict, errors.New("Already registered.")
	}

	var key *datastore.Key
	if key, err = company.Save(ctx); err != nil {
		return http.StatusInternalServerError, err
	}

	// TODO(flowlo): Respond with HTTP 201 and include a
	// location header and caching information.

	company.Write(w, key)
	return http.StatusOK, nil
}
