package util

import (
	"errors"
)

// CheckNick should be used to check whether a nick can be accepted.
func CheckNick(nick string) error {
	if nick == "" {
		return errors.New("nick is zero")
	}

	bytes := []byte(nick)

	if bytes[0] == 45 {
		return errors.New("nick starts with hyphen")
	}
	if bytes[len(bytes)-1] == 45 {
		return errors.New("nick ends with hyphen")
	}

	var prev byte
	for _, b := range bytes {
		// no number              no uppercase letter    no lowercase letter     no hyphen
		if !(b > 47 && b < 58) && !(b > 64 && b < 91) && !(b > 96 && b < 123) && b != 45 {
			return errors.New("nick has invalid symbol")
		}
		if prev == 45 && b == 45 {
			return errors.New("nick has two consecutive hyphens")
		}
		prev = b
	}
	return nil
}
