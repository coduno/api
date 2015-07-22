package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/app/models"
	"github.com/coduno/app/util"
)

// CompanyLoginInfo is the login info for a company
type CompanyLoginInfo struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

// CompanyLogin starts a session for a company
func CompanyLogin(w http.ResponseWriter, r *http.Request, c context.Context) (createSession bool) {
	createSession = false
	if !util.CheckMethod(w, r, "POST") {
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Error reading: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var companyLogin CompanyLoginInfo
	err = json.Unmarshal(body, &companyLogin)

	if err != nil {
		http.Error(w, "Cannot unmarshal: "+err.Error(), http.StatusInternalServerError)
		return
	}
	q := datastore.NewQuery("companies").Filter("Email = ", companyLogin.Email).Limit(1)
	var companies []models.Company
	keys, err := q.GetAll(c, &companies)
	if err != nil {
		http.Error(w, "Datastore error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(companies) != 1 {
		http.Error(w, "You are unauthorized to login!", http.StatusUnauthorized)
		return
	}


	company := companies[0]
	if !util.CheckPassword(company.Password, companyLogin.Password){
		http.Error(w, "You are unauthorized to login!", http.StatusUnauthorized)
		return
	}
	company.EntityID = keys[0].Encode()

	toSend := make(map[string]interface{})
	toSend["company"] = models.Company{Name: company.Name, Email: company.Email, EntityID: company.EntityID}
	json, err := json.Marshal(toSend)
	if err != nil {
		http.Error(w, "Json marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(json))
	return true
}

func CreateCompany(w http.ResponseWriter, r *http.Request, c context.Context){
	if !util.CheckMethod(w, r, "POST") {
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Error reading: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var company models.Company
	err = json.Unmarshal(body, &company)

	if err != nil {
		http.Error(w, "Cannot unmarshal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	q := datastore.NewQuery("companies").Filter("Email = ", company.Email).Limit(1)
	var companies []models.Company
	_, err = q.GetAll(c, &companies)
	if err != nil {
		http.Error(w, "Datastore error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(companies) > 0 {
		toSend := make(map[string]interface{})
		toSend["error"] = "email already exists"
		json, _ := json.Marshal(toSend)
		w.Write([]byte(json))
		return
	}

	company = company.SaveCompany(c)

	json, err := json.Marshal(company)
	w.Write([]byte(json))
	return
}
