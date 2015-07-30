package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/engine/model"
	"github.com/coduno/engine/util"
	"github.com/coduno/engine/util/password"
)

// LoginInfo is the login info for a company
type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CompanyLogin starts a session for a company
func CompanyLogin(c context.Context, w http.ResponseWriter, r *http.Request) {
	if !util.CheckMethod(w, r, "POST") {
		return
	}

	var err error

	var body []byte
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var loginInfo LoginInfo
	if err = json.Unmarshal(body, &loginInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	q := model.NewQueryForUser().
		Filter("Address = ", loginInfo.Email).
		Limit(1)

	var users model.Users
	keys, err := q.GetAll(c, &users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(users) != 1 {
		// NOTE: Do not leak len(users) here.
		http.Error(w, "permission denied", http.StatusUnauthorized)
		return
	}

	user := users[0]
	key := keys[0]

	if err = password.Check([]byte(loginInfo.Password), user.HashedPassword); err != nil {
		// NOTE: Do not leak err here.
		http.Error(w, "permission denied", http.StatusUnauthorized)
		return
	}

	user.Write(w, key)
}

func PostCompany(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if !util.CheckMethod(w, r, "POST") {
		return
	}
	var err error

	var company model.Company
	if err = json.NewDecoder(r.Body).Decode(&company); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	q := model.NewQueryForCompany().
		Filter("Address = ", company.Address.Address).
		Limit(1)

	var companies model.Companys
	if _, err = q.GetAll(ctx, &companies); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(companies) > 0 {
		body, _ := json.Marshal(map[string]string{
			"error": "Company already exists",
		})
		w.Write(body)
		return
	}

	var key *datastore.Key
	if key, err = company.Save(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO(flowlo): Respond with HTTP 201 and include a
	// location header and caching information.

	company.Write(w, key)
}
