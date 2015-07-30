package model

import "net/mail"

//go:generate generator -c "Profile"

// User is anybody interacting with our systems. It will
// ultimately refer to who accessed Coduno (or on whose
// behalf).
//
// Saved in Datastore, User will be optionally a child
// entity to Company, so keys pointing to a User
// can be used to obtain the company they work.
type User struct {
	// Encapsulates Name (combined first and last name,
	// however the user likes) and an e-mail address.
	//
	// Datastore will split this into two properties
	// called Name and Address, where Address must be
	// guaranteed to be unique.
	mail.Address

	// Unique name for this user, like analogous to @flowlo
	// on GitHub/Twitter/... Mandatory and unique.
	Nick string

	// Hashed and salted password to be accessed by
	// corresponding helpers in util.
	// See https://godoc.org/golang.org/x/crypto/bcrypt
	HashedPassword []byte `json:"-"`
}
