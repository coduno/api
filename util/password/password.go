package password

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// MinLength is the minimum length for any password to be
// generated, hashed, or checked by functions in this package.
const MinLength = 12

// MaxLength is the maximum length for any password to be
// generated, hashed, or checked by functions in this package.
const MaxLength = 2048

// DefaultLength is the length chosen by Generate in case the
// passed length is smaller than MinLength.
const DefaultLength = 16

// InvalidPasswordLengthError is returned by Hash and Check
// in case the given plaintext password is too short or too
// long.
// Generate is guaranteed to produce passwords that adhere
// to those limits.
type InvalidPasswordLengthError int

func (ipl InvalidPasswordLengthError) Error() string {
	return fmt.Sprintf("app/util/password: length %d is outside allowed range (%d,%d)", int(ipl), int(MinLength), int(MaxLength))
}

// Generate produces a cryptographically strong random
// slice of bytes.
// If n is lower than MinLength it is set to
// DefaultLength.
func Generate(n int) (password []byte, err error) {
	if n < MinLength {
		n = DefaultLength
	}
	password = make([]byte, n)
	_, err = rand.Read(password)
	return password, err
}

// Hash computes a hash from the given password.
func Hash(password []byte) (hash []byte, err error) {
	if err = checkLen(password); err != nil {
		return
	}

	return bcrypt.GenerateFromPassword(password, 0)
}

// Check compares the given password (plaintext) against
// some hash. It will return nil if password and hash
// match.
func Check(password, hash []byte) (err error) {
	if err = checkLen(password); err != nil {
		return
	}

	return bcrypt.CompareHashAndPassword(hash, password)
}

func checkLen(password []byte) error {
	if len(password) < MinLength || len(password) > MaxLength {
		return InvalidPasswordLengthError(len(password))
	}

	return nil
}
