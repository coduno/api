// Package passenger revolves around the User associated with the
// request in flight.
//
// It establishes this link via the HTTP authorization header,
// which can be used to authenticate with basic authorization or
// token based authorization.
package passenger

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/subtle"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	mathrand "math/rand"
	"net/http"
	"strings"
	"time"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util/password"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"

	// Indirect import of SHA2, will be called by crypto.Hash.New()
	_ "crypto/sha256"
	_ "crypto/sha512"

	// Indirect import of SHA3, will be called by crypto.Hash.New()
	_ "golang.org/x/crypto/sha3"
)

const maxValidity = time.Hour * 24 * 365 * 2
const minValidity = time.Minute * 30
const defaultValidity = time.Hour * 24 * 14

const defaultHash = crypto.SHA3_256

type key int64

var passengerKey key

// ErrNoAuthHeader signals that the request did not carry an Authorization
// header.
var ErrNoAuthHeader = errors.New("passenger: no authorization header present")

// ErrUnkAuthHeader is returned if the request contained an Authorization
// header, but it could not be parsed. The two accepted authorization
// methods are "Basic" and "Token".
var ErrUnkAuthHeader = errors.New("passenger: cannot interpret authorization header")

// ErrDigestMismatch is returned if the token provided by the user did not
// result in the same digest value after hashing as the one stored on the
// server.
var ErrDigestMismatch = errors.New("passenger: digest mismatch")

// ErrTokenNotAssociated can be returned if this package encounters invalid
// data during authentication. An AccessToken must be child of a User.
var ErrTokenNotAssociated = errors.New("passenger: token not associated to any user")

// ErrTokenExpired is returned if the token is not valid anymore.
var ErrTokenExpired = errors.New("passenger: token expired")

// ErrTokenNotMatchingUser is returned if basic auth with token password
// was attempted, and the token could be found but does not match the user.
type ErrTokenNotMatchingUser struct {
	Parent, Actual *datastore.Key
}

func (e ErrTokenNotMatchingUser) Error() string {
	return fmt.Sprintf("passenger: parent of token %+v does not match with actual user %+v", e.Parent, e.Actual)
}

func init() {
	passengerKey = key(mathrand.Int63())
}

// Passenger holds the currently authenticated user
// together with the access token (if relevant).
type Passenger struct {
	UserKey     *datastore.Key
	AccessToken model.AccessToken
}

// HasScope can be used to check whether this Passenger was
// granted access to a given scope.
// Note that a user authenticated via username and password
// (not via access token) will have access to all scopes by
// default.
func (p *Passenger) HasScope(scope string) (has bool) {
	for _, grantedScope := range p.AccessToken.Scopes {
		if scope == grantedScope {
			return true
		}
	}
	return
}

// Save will persist the passenger to Datastore and send it to Memcache.
func (p *Passenger) Save(ctx context.Context) (*datastore.Key, error) {
	now := time.Now()

	key, err := p.AccessToken.SaveWithParent(ctx, p.UserKey)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err = gob.NewEncoder(buf).Encode(p); err != nil {
		return nil, err
	}

	item := &memcache.Item{
		Key:        key.Encode(),
		Value:      buf.Bytes(),
		Expiration: p.AccessToken.Expiry.Sub(now) + 10*time.Second,
	}

	if err = memcache.Set(ctx, item); err != nil {
		return nil, err
	}

	return key, nil
}

func (p *Passenger) check(raw []byte) error {
	digest := crypto.Hash(p.AccessToken.Hash).New().Sum(raw)

	if subtle.ConstantTimeCompare(digest, p.AccessToken.Digest) != 1 {
		return ErrDigestMismatch
	}

	if p.AccessToken.Expiry.Before(time.Now()) {
		return ErrTokenExpired
	}

	return nil
}

// IssueToken creates a new Token for the authenticated user. Callers should
// prefill the accessToken with whatever values they like and leave zero values
// to be set.
// The generated token will also be persisted and can be handed to the client
// with no more handling.
func (p *Passenger) IssueToken(ctx context.Context, token *model.AccessToken) (string, error) {
	now := time.Now()

	if token.Expiry == (time.Time{}) {
		token.Expiry = now.Add(defaultValidity)
	} else {
		if token.Expiry.Before(now.Add(minValidity)) {
			return "", fmt.Errorf("passenger: token must be valid at least %s", minValidity)
		}

		if token.Expiry.After(now.Add(maxValidity)) {
			return "", fmt.Errorf("passenger: token must be valid at most %s", maxValidity)
		}
	}

	// TODO(flowlo): This will reject all scopes for now, as we are not using them.
	// As soon as we introduce scopes, this check must be rewritten accordingly.
	if len(token.Scopes) > 0 {
		return "", fmt.Errorf("passenger: unknown scopes: %s", strings.Join(token.Scopes, ", "))
	}

	if len(token.Description) > 512 || len(token.Description) < 4 {
		return "", fmt.Errorf("passenger: description has bad len %d", len(token.Description))
	}

	raw := make([]byte, 16)

	if _, err := rand.Read(raw); err != nil {
		return "", err
	}

	if token.Hash == 0 {
		token.Hash = int(defaultHash)
	}
	token.Digest = crypto.Hash(token.Hash).New().Sum(raw)

	clone := Passenger{
		UserKey:     p.UserKey,
		AccessToken: *token,
	}

	key, err := clone.Save(ctx)
	if err != nil {
		return "", err
	}

	return encodeToken(key, raw)
}

// FromAccessToken tries do identify a Passenger by the access token he gave us.
// It will look up the AccessToken and consequently the corresponding User.
func FromAccessToken(ctx context.Context, accessToken string) (*Passenger, error) {
	key, raw, err := decodeToken(accessToken)
	if err != nil {
		return nil, err
	}

	var p *Passenger
	if p, err = fromCache(ctx, key); err != nil {
		if p, err = fromDatastore(ctx, key); err != nil {
			return nil, err
		}
	}

	return p, p.check(raw)
}

func fromCache(ctx context.Context, key *datastore.Key) (p *Passenger, err error) {
	item, err := memcache.Get(ctx, key.Encode())
	if err != nil {
		return nil, err
	}
	p = new(Passenger)
	err = gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&p)
	return
}

func fromDatastore(ctx context.Context, key *datastore.Key) (p *Passenger, err error) {
	p = new(Passenger)
	if err = datastore.Get(ctx, key, p.AccessToken); err != nil {
		return
	}
	if p.UserKey = key.Parent(); p.UserKey == nil {
		return nil, ErrTokenNotAssociated
	}
	return
}

// FromBasicAuth tries do identify a Passenger by the access token he gave us.
// It will look up the the user by username and try to match password.
func FromBasicAuth(ctx context.Context, username, pw string) (p *Passenger, err error) {
	p = new(Passenger)
	var user model.User
	p.UserKey, err = model.NewQueryForUser().
		Filter("Nick=", username).
		Limit(1).
		Run(ctx).
		Next(&user)

	if err != nil {
		return
	}
	err = password.Check([]byte(pw), user.HashedPassword)

	// TODO(flowlo): Depending on bcrypt is very fragile. We
	// should encapsulate that.
	if err == bcrypt.ErrMismatchedHashAndPassword {
		userKey := p.UserKey
		p, err = FromAccessToken(ctx, pw)
		if err != nil {
			return
		}
		if !p.UserKey.Equal(userKey) {
			return nil, ErrTokenNotMatchingUser{Parent: p.UserKey, Actual: userKey}
		}
	}
	return
}

// FromRequest inspects the HTTP Authorization header of the given request
// and tries to identify a passenger.
func FromRequest(ctx context.Context, r *http.Request) (p *Passenger, err error) {
	auth := ""
	if auth = r.Header.Get("Authorization"); auth == "" {
		return nil, ErrNoAuthHeader
	}

	if strings.HasPrefix(auth, "Token ") {
		return FromAccessToken(ctx, auth[6:])
	}

	username, password, ok := "", "", false
	if username, password, ok = r.BasicAuth(); !ok {
		return nil, ErrUnkAuthHeader
	}

	p, err = FromBasicAuth(ctx, username, password)
	return
}

// FromContext returns the Passenger value stored in ctx if any.
func FromContext(ctx context.Context) (p *Passenger, ok bool) {
	p, ok = ctx.Value(passengerKey).(*Passenger)
	return
}

// NewContext returns a new Context that carries value p.
func NewContext(ctx context.Context, p *Passenger) context.Context {
	return context.WithValue(ctx, passengerKey, p)
}

// NewContextFromRequest wraps FromRequest and NewContext.
func NewContextFromRequest(ctx context.Context, r *http.Request) (context.Context, error) {
	p, err := FromRequest(ctx, r)
	if err != nil {
		return ctx, err
	}
	return NewContext(ctx, p), nil
}

// DecodeToken will take a token as sent by a client and translate it into a
// key to look up the full token on the server side and the raw secret.
func decodeToken(token string) (*datastore.Key, []byte, error) {
	b, err := hex.DecodeString(token)
	if err != nil {
		return nil, nil, err
	}

	dec := gob.NewDecoder(bytes.NewReader(b))

	var key *datastore.Key
	if err := dec.Decode(&key); err != nil {
		return nil, nil, err
	}
	var raw []byte
	if err := dec.Decode(&raw); err != nil {
		return nil, nil, err
	}
	return key, raw, nil
}

// EncodeToken translates the key and raw secret of a newly generated token to
// a form suitable for the client.
func encodeToken(key *datastore.Key, raw []byte) (string, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(key); err != nil {
		return "", err
	}
	if err := enc.Encode(raw); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}
