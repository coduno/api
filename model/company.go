package model

import "net/mail"

//go:generate generator -c "Challenge"

// Company contains the data related to a company.
//
// TODO(flowlo, victorbalan): In the future, the company
// may point at Users to enable role based authentication.
type Company struct {
	mail.Address `datastore:",index",json:",omitempty"`

	// Unique name for this user, like analogous to @google
	// on GitHub/Twitter/...
	Nick string `datastore:",index",json:",omitempty"`
}

func (c Company) IsValid() bool {
	return true
}
