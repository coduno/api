package passenger

import (
	"encoding/hex"
	"errors"
	"math/rand"
	"net/http"
	"strings"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/app/models"
	"github.com/coduno/app/util/password"

	"golang.org/x/net/context"
)

type key int64

var passengerKey key

var ErrNoAuthHeader = errors.New("passenger: no authorization header present")
var ErrUnkAuthHeader = errors.New("passenger: cannot interpret authorization header")

// Passenger holds the currently authenticated user
// together with the access token (if relevant).
type Passenger struct {
	User           *models.User
	UserKey        *datastore.Key
	AccessToken    *models.AccessToken
	AccessTokenKey *datastore.Key
}

// HasScope can be used to check whether this Passenger was
// granted access to a given scope.
// Note that a user authenticated via username and password
// (not via access token) will have access to all scopes by
// default.
func (p *Passenger) HasScope(scope string) (has bool) {
	if p.AccessToken == nil {
		return true
	}
	for _, grantedScope := range p.AccessToken.Scopes {
		if scope == grantedScope {
			return true
		}
	}
	return
}

// FromAccessToken tries do identify a Passenger by the access token he gave us.
// It will look up the AccessToken and consequently the corresponding User.
func FromAccessToken(ctx context.Context, accessToken []byte) (p *Passenger, err error) {
	p.AccessTokenKey, err = datastore.NewQuery(models.AccessTokenKind).
		Filter("Value=", accessToken).
		Limit(1).
		Run(ctx).
		Next(p.AccessToken)

	if err != nil {
		return
	}
	if p.UserKey = p.AccessTokenKey.Parent(); p.UserKey == nil {
		return
	}
	err = datastore.Get(ctx, p.UserKey, p.User)
	return
}

// FromBasicAuth tries do identify a Passenger by the access token he gave us.
// It will look up the the user by username and try to match password.
func FromBasicAuth(ctx context.Context, username, pw string) (p *Passenger, err error) {
	p.UserKey, err = datastore.NewQuery(models.UserKind).
		Filter("Username=", username).
		Limit(1).
		Run(ctx).
		Next(p.User)

	err = password.Check([]byte(pw), p.User.HashedPassword)
	return
}

// FromRequest inspects the HTTP Authorization header of the given request
// and tries to identify a passenger.
func FromRequest(ctx context.Context, r *http.Request) (p *Passenger, err error) {
	if username, password, ok := r.BasicAuth(); ok {
		return FromBasicAuth(ctx, username, password)
	}

	auth := ""
	if auth = r.Header.Get("Authorization"); auth == "" {
		return nil, ErrNoAuthHeader
	}

	if !strings.HasPrefix(auth, "Token ") {
		return nil, ErrUnkAuthHeader
	}

	var token []byte
	if _, err := hex.Decode(token, []byte(auth[6:])); err != nil {
		return nil, ErrUnkAuthHeader
	}

	return FromAccessToken(ctx, token)
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
		return nil, err
	}
	return NewContext(ctx, p), nil
}

func init() {
	passengerKey = key(rand.Int63())
}
