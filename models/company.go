package models

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// CompanyKind is the kind used to store companies in
// Datastore.
const CompanyKind = "companies"

// Company contains the data related to a company.
type Company struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	HashedPassword []byte `json:"-"`
}

// Save puts this company in the Datastore.
func (c *Company) Save(ctx context.Context) (key *datastore.Key, err error) {
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, CompanyKind, nil), c)
}
