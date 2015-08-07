package model

import "time"

//go:generate generator

// AccessToken encapsulates a string that be used to
// authenticate a User.
type AccessToken struct {
	// Corresponds to the crypto.Hash that was used to hash the value
	// of this AccessToken.
	// It is an int and not a crypto.Hash because the base type of
	// crypto.Hash is uint and unsigned types cannot be stored by
	// Datastore.
	// A conversion is needed at runtime:
	//
	//	hash := crypto.Hash(accessToken.Hash)
	//
	Hash int `datastore:",noindex"`

	// The digest of this AccessToken after hashing with above hash.
	Digest []byte `datastore:",noindex"`

	// If a User is authenticated using an AccessToken, authorization
	// can be granted to only a subset of possible actions. This slice
	// acts as a filter and should list allowed scopes, i.e. permissions.
	Scopes []string `datastore:",noindex"`

	// Arbitrary string describing the use of this token. It can be
	// automatically generated or set by the user.
	Description string `datastore:",noindex"`

	// Time of creation.
	Creation time.Time `datastore:",noindex"`

	// If an AccessToken is seen after Expiry, it is to be deleted.
	Expiry time.Time `datastore:",index"`

	// Address of the client that created this AccessToken.
	RemoteAddr string `datastore:",index"`
}
