package passenger

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util/password"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type key int64

var passengerKey key

var ErrNoAuthHeader = errors.New("passenger: no authorization header present")
var ErrUnkAuthHeader = errors.New("passenger: cannot interpret authorization header")
var ErrTokenNotAssociated = errors.New("")

type ErrTokenNotMatchingUser struct {
	Parent, Actual *datastore.Key
}

func (e ErrTokenNotMatchingUser) Error() string {
	return fmt.Sprintf("passenger: parent of token %+v does not match with actual user %+v", e.Parent, e.Actual)
}

// Passenger holds the currently authenticated user
// together with the access token (if relevant).
type Passenger struct {
	User           model.User
	UserKey        *datastore.Key
	AccessToken    model.AccessToken
	AccessTokenKey *datastore.Key
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

// FromAccessToken tries do identify a Passenger by the access token he gave us.
// It will look up the AccessToken and consequently the corresponding User.
func FromAccessToken(ctx context.Context, accessToken string) (p *Passenger, err error) {
	p = new(Passenger)
	p.AccessTokenKey, err = model.NewQueryForAccessToken().
		Filter("Value=", accessToken).
		Limit(1).
		Run(ctx).
		Next(&p.AccessToken)

	if err != nil {
		return
	}
	if p.UserKey = p.AccessTokenKey.Parent(); p.UserKey == nil {
		return nil, ErrTokenNotAssociated
	}
	// TODO(flowlo): Make this independent of Datastore import.
	err = datastore.Get(ctx, p.UserKey, &p.User)
	return
}

// FromBasicAuth tries do identify a Passenger by the access token he gave us.
// It will look up the the user by username and try to match password.
func FromBasicAuth(ctx context.Context, username, pw string) (p *Passenger, err error) {
	p = new(Passenger)
	p.UserKey, err = model.NewQueryForUser().
		Filter("Nick=", username).
		Limit(1).
		Run(ctx).
		Next(&p.User)

	if err != nil {
		return
	}
	err = password.Check([]byte(pw), p.User.HashedPassword)

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

func init() {
	passengerKey = key(rand.Int63())
}
