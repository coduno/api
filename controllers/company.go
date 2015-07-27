package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/app/models"
	"github.com/coduno/app/util"
	"github.com/coduno/app/util/password"
)

// CompanyLoginInfo is the login info for a company
type CompanyLoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CompanyLogin starts a session for a company
func CompanyLogin(w http.ResponseWriter, r *http.Request, c context.Context) (createSession bool) {
	if !util.CheckMethod(w, r, "POST") {
		return
	}

	var err error

	var body []byte
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var companyLogin CompanyLoginInfo
	if err = json.Unmarshal(body, &companyLogin); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	q := datastore.NewQuery(models.CompanyKind).
		Filter("Email = ", companyLogin.Email).
		Limit(1)

	var companies []models.Company
	keys, err := q.GetAll(c, &companies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(companies) != 1 {
		// NOTE: Do not leak len(companies) here.
		http.Error(w, "permission denied", http.StatusUnauthorized)
		return
	}

	company := companies[0]
	key := keys[0]

	if err = password.Check([]byte(companyLogin.Password), company.HashedPassword); err != nil {
		// NOTE: Do not leak err here.
		http.Error(w, "permission denied", http.StatusUnauthorized)
		return
	}

	util.WriteEntity(w, key, company)
	return true
}

func CreateCompany(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if !util.CheckMethod(w, r, "POST") {
		return
	}

	var err error

	var body []byte
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var company models.Company
	if err = json.Unmarshal(body, &company); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	q := datastore.NewQuery(models.CompanyKind).
		Filter("Email = ", company.Email).
		Limit(1)

	var companies []models.Company
	if _, err = q.GetAll(ctx, &companies); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(companies) > 0 {
		body, _ := json.Marshal(map[string]string{
			"error": "email already exists",
		})
		w.Write(body)
		return
	}

	var pw []byte
	if pw, err = password.Generate(0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var hpw []byte
	if hpw, err = password.Hash(pw); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	company.HashedPassword = hpw

	var key *datastore.Key
	if key, err = company.Save(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO(flowlo): Respond with HTTP 201 and include a
	// location header and caching information.

	util.WriteEntity(w, key, company)
	return
}
