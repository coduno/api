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
	Company string `json:"company"`
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
	q := datastore.NewQuery("companies").Filter("Name = ", "Catalysts").Limit(1)
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
	company.EntityID = keys[0].Encode()

	toSend := make(map[string]interface{})
	toSend["company"] = company
	json, err := json.Marshal(toSend)
	if err != nil {
		http.Error(w, "Json marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(json))
	return true
}
