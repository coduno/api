package model

import (
	"net/mail"

	"google.golang.org/appengine/datastore"
)

//go:generate generator -c "Profile"

// User is anybody interacting with our systems. It will
// ultimately refer to who accessed Coduno (or on whose
// behalf).
type User struct {
	// Encapsulates Name (combined first and last name,
	// however the user likes) and an e-mail address.
	//
	// Datastore will split this into two properties
	// called Name and Address, where Address must be
	// guaranteed to be unique.
	mail.Address `datastore:",index",json:",omitempty"`

	// Unique name for this user, like analogous to @flowlo
	// on GitHub/Twitter/...
	Nick string `datastore:",index",json:",omitempty"`

	// Points to the company a user works for, if any.
	Company *datastore.Key `datastore:",index",json:",omitempty"`

	// Hashed and salted password to be accessed by
	// corresponding helpers in util.
	// See https://godoc.org/golang.org/x/crypto/bcrypt
	HashedPassword []byte `datastore:",noindex" json:"-"`
}
