package model

import "net/mail"

//go:generate generator -c "Challenge,User"

// Company contains the data related to a company.
//
// TODO(flowlo, victorbalan): In the future, the company
// will point at Users to enable role based authentication.
type Company struct {
	mail.Address
}
