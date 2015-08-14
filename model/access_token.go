package model

import "time"

//go:generate generator

// AccessToken encapsulates a digest of a secret that can be used to
// authenticate a User. The entity only holds a digest to prevent impersonation
// of a user in case it is leaked.
//
// When requests arrive, the correct AccessToken is queried by key, therefore
// these entities use indexes vary sparingly.
//
// AccessTokens reside in an entity group rooted at a User. As AccessTokens
// cannot be altered by the user, writes are only done at creation and deletion
// and are therefore neglegible.
type AccessToken struct {
	// Corresponds to the crypto.Hash that was used to hash the value
	// of this AccessToken.
	// It is an int and not a crypto.Hash because the base type of
	// crypto.Hash is uint (go1) and unsigned types cannot be stored by
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

	// If an AccessToken is seen after Expiry, it is to be deleted. This property
	// is indexed to enable for grabage collection of expired AccessTokens.
	Expiry time.Time `datastore:",index"`

	// Address of the client that created this AccessToken.
	RemoteAddr string `datastore:",noindex"`
}
