package controllers

import (
	"encoding/base64"
	"strings"
)

// EncodeToken creates a token from the username and the fingerprintId. This will be needed
// to authenticate the client.
func EncodeToken(username string, fingerprintID string) string {
	var token = username + ":" + fingerprintID
	encodedToken := base64.StdEncoding.EncodeToString([]byte(token))

	return encodedToken
}

// DecodeToken decodes the token received from the client and returns the parsed username
// and fingerprintId
func DecodeToken(token string) (string, string, error) {
	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", "", err
	}
	splitToken := strings.Split(string(decodedToken), ":")
	return splitToken[0], splitToken[1], nil
}
