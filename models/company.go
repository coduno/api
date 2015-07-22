package models

import(
	"golang.org/x/net/context"
	"github.com/coduno/app/util"
	"google.golang.org/appengine/datastore"
)

// Company contains the data related to a company
type Company struct {
	EntityID string `datastore:"-" json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SaveCompany(company Company, ctx context.Context) Company {
	  password := util.RandomPassword()
		//this random password will be mailed to the company at company.Email
		company.Password = util.GeneratePassword(password)
		datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "companies", nil), &company)
		return company
}
