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
	"encoding/binary"
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

const tokenLength = 8

var kinds = [2]string{model.UserKind, model.TokenKind}

var endianness = binary.LittleEndian

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
// data during authentication. A Token must be child of a User.
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
// together with the token (if relevant).
type Passenger struct {
	User  *datastore.Key
	Token *model.Token
}

// HasScope can be used to check whether this Passenger was
// granted access to a given scope.
// Note that a user authenticated via username and password
// (not via token) will have access to all scopes by
// default.
func (p *Passenger) HasScope(scope string) (has bool) {
	for _, grantedScope := range p.Token.Scopes {
		if scope == grantedScope {
			return true
		}
	}
	return
}

// Save will persist the passenger to Datastore and send it to Memcache.
func (p *Passenger) Save(ctx context.Context) (*datastore.Key, error) {
	now := time.Now()

	key, err := p.Token.PutWithParent(ctx, p.User)
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
		Expiration: p.Token.Expiry.Sub(now) + 10*time.Second,
	}

	if err = memcache.Set(ctx, item); err != nil {
		return nil, err
	}

	return key, nil
}

func (p *Passenger) check(raw []byte) error {
	digest := crypto.Hash(p.Token.Hash).New().Sum(raw)

	if subtle.ConstantTimeCompare(digest, p.Token.Digest) != 1 {
		return ErrDigestMismatch
	}

	if p.Token.Expiry.Before(time.Now()) {
		return ErrTokenExpired
	}

	return nil
}

// IssueToken creates a new Token for the authenticated user. Callers should
// prefill the Token with whatever values they like and leave zero values
// to be set.
// The generated token will also be persisted and can be handed to the client
// with no more handling.
func (p *Passenger) IssueToken(ctx context.Context, token *model.Token) (string, error) {
	now := time.Now()

	token.Creation = now

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

	var raw [tokenLength]byte

	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}

	if token.Hash == 0 {
		token.Hash = int(defaultHash)
	}
	token.Digest = crypto.Hash(token.Hash).New().Sum(raw[:])

	clone := Passenger{
		User:  p.User,
		Token: token,
	}

	key, err := clone.Save(ctx)
	if err != nil {
		return "", err
	}

	return encodeToken(key, &raw)
}

// FromToken tries do identify a Passenger by the token he gave us.
// It will look up the Token and consequently the corresponding User.
func FromToken(ctx context.Context, Token string) (*Passenger, error) {
	key, raw, err := decodeToken(ctx, Token)
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

func fromCache(ctx context.Context, key *datastore.Key) (*Passenger, error) {
	item, err := memcache.Get(ctx, key.Encode())
	if err != nil {
		return nil, err
	}
	p := new(Passenger)
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&p); err != nil {
		return nil, err
	}
	return p, nil
}

func fromDatastore(ctx context.Context, key *datastore.Key) (*Passenger, error) {
	p := &Passenger{
		Token: &model.Token{},
	}
	if err := datastore.Get(ctx, key, p.Token); err != nil {
		return nil, err
	}
	if p.User = key.Parent(); p.User == nil {
		return nil, ErrTokenNotAssociated
	}
	return p, nil
}

// FromBasicAuth tries do identify a Passenger by the token he gave us.
// It will look up the the user by username and try to match password.
func FromBasicAuth(ctx context.Context, username, pw string) (p *Passenger, err error) {
	p = new(Passenger)
	var user model.User
	p.User, err = model.NewQueryForUser().
		Filter("Nick =", username).
		Limit(1).
		Run(ctx).
		Next(&user)

	if err != nil {
		return nil, err
	}
	err = password.Check([]byte(pw), user.HashedPassword)

	// TODO(flowlo): Depending on bcrypt is very fragile. We
	// should encapsulate that.
	if err == bcrypt.ErrMismatchedHashAndPassword {
		userKey := p.User
		p, err = FromToken(ctx, pw)
		if err != nil {
			return
		}
		if !p.User.Equal(userKey) {
			return nil, ErrTokenNotMatchingUser{Parent: p.User, Actual: userKey}
		}
	}
	return
}

// FromRequest inspects the HTTP Authorization header of the given request
// and tries to identify a passenger.
func FromRequest(ctx context.Context, r *http.Request) (*Passenger, error) {
	auth := ""
	if auth = r.Header.Get("Authorization"); auth == "" {
		return nil, ErrNoAuthHeader
	}

	if strings.HasPrefix(auth, "Token ") {
		return FromToken(ctx, auth[6:])
	}

	username, password, ok := "", "", false
	if username, password, ok = r.BasicAuth(); !ok {
		return nil, ErrUnkAuthHeader
	}

	return FromBasicAuth(ctx, username, password)
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
func decodeToken(ctx context.Context, token string) (*datastore.Key, []byte, error) {
	b, err := hex.DecodeString(token)
	if err != nil {
		return nil, nil, err
	}

	if len(b) != len(kinds)*8+tokenLength {
		return nil, nil, errors.New("token length mismatch")
	}

	var intIDs [len(kinds)]int64
	for i := range intIDs {
		intID, n := binary.Varint(b[i*8 : (i+1)*8])
		if n < 8 {
			return nil, nil, errors.New("varint read mismatch")
		}
		intIDs[len(intIDs)-1-i] = intID
	}

	var key *datastore.Key
	for i := range intIDs {
		if intIDs[i] == 0 {
			continue
		}
		key = datastore.NewKey(ctx, kinds[i], "", intIDs[i], key)
		if key == nil {
			return nil, nil, errors.New("could not construct key from token")
		}
	}

	return key, b[len(kinds)*8:], nil
}

// EncodeToken translates the key and raw secret of a newly generated token to
// a form suitable for the client.
func encodeToken(key *datastore.Key, raw *[tokenLength]byte) (string, error) {
	// Buffer size will be 8 (size of an int64) times the number of keys
	// in the hirarchy plus the length of the raw token itself.
	var b [len(kinds)*8 + tokenLength]byte

	for i := range kinds {
		if n := binary.PutVarint(b[i*8:(i+1)*8], key.IntID()); n < 8 {
			return "", errors.New("short write when encoding token")
		}
		if key != nil {
			key = key.Parent()
		}
	}

	copy(b[len(kinds)*8:len(kinds)*8+tokenLength], raw[:])

	return hex.EncodeToString(b[:]), nil
}
