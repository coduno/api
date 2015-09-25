package controllers

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"testing"

	"github.com/coduno/api/model"

	"google.golang.org/appengine/datastore"
)

func TestPostCompany(t *testing.T) {
	company := model.Company{
		Address: mail.Address{
			Name:    "companyName",
			Address: "email@example.com",
		},
	}

	r, err := http.NewRequest("POST", "/companies", requestBody(t, company))
	if err != nil {
		t.Fatal(err)
	}
	loginAsCompanyUser(r)
	rr := record(t, r)
	testRequestStatus(t, rr, 200, "Should be okay.")

	// Check the entity the client recieves was indeed saved in the datastore.
	var companyResponse = struct {
		Name, Address, Nick, Key string
	}{}

	if err = json.NewDecoder(rr.Body).Decode(&companyResponse); err != nil {
		t.Fatal(err)
	}

	var key *datastore.Key
	key, err = datastore.DecodeKey(companyResponse.Key)
	if err != nil {
		t.Fatal(err)
	}

	ctx := backgroundContext()

	var savedCompany model.Company

	err = datastore.Get(ctx, key, &savedCompany)
	if err != nil {
		t.Fatal(err)
	}

	if savedCompany.Address != company.Address {
		t.Fatal("Saved company is wrong")
	}

	datastore.Delete(ctx, key)
}
