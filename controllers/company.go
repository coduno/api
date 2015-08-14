package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/api/model"
)

// PostCompany creates a new company after validating by key.
func PostCompany(ctx context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, nil
	}

	var company model.Company
	if err = json.NewDecoder(r.Body).Decode(&company); err != nil {
		return http.StatusBadRequest, err
	}

	var companies model.Companys
	_, err = model.NewQueryForCompany().
		Filter("Address = ", company.Address.Address).
		Limit(1).
		GetAll(ctx, &companies)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if len(companies) > 0 {
		return http.StatusConflict, errors.New("Already registered.")
	}

	var key *datastore.Key
	if key, err = company.Save(ctx, nil); err != nil {
		return http.StatusInternalServerError, err
	}

	// TODO(flowlo): Respond with HTTP 201 and include a
	// location header and caching information.

	json.NewEncoder(w).Encode(company.Key(key))
	return http.StatusOK, nil
}
