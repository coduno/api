package controllers

import (
	"encoding/json"
	"net/mail"
	"testing"

	"github.com/coduno/api/model"
	"github.com/coduno/api/tests/testUtils"

	"google.golang.org/appengine/datastore"
)

func TestPostCompany(t *testing.T) {
	company := model.Company{
		Address: mail.Address{
			Name:    "companyName",
			Address: "email@example.com",
		},
	}

	req, w := requestAndResponse(t, "POST", "/companies", company)
	ctx = testUtils.LoginAsCompanyUser(t, ctx, req)
	r.ServeHTTP(w, req)
	testRequestStatusAndError(t, 200, nil)

	// Check the entity the client recieves was indeed saved in the datastore.
	var companyResponse = struct {
		Name,
		Address,
		Nick,
		Key string
	}{}
	var err error

	if err = json.NewDecoder(w.Body).Decode(&companyResponse); err != nil {
		t.Fatal(err)
	}

	var key *datastore.Key
	key, err = datastore.DecodeKey(companyResponse.Key)
	if err != nil {
		t.Fatal(err)
	}

	var savedCompany model.Company

	err = datastore.Get(ctx, key, &savedCompany)
	if err != nil {
		t.Fatal(err)
	}

	if savedCompany.Address != company.Address {
		t.Fatal("Saved company is wrong")
	}

	datastore.Delete(ctx, key)
	ctx = testUtils.Logout(ctx)
}
