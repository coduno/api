package controllers

import (
	"net/http"
	"testing"
)

func TestGetChallengeByKey(t *testing.T) {
	key := challengeKey.Encode()
	r, err := http.NewRequest("GET", "/challenges/"+key, nil)
	if err != nil {
		t.Fatal(err)
	}
	loginAsCompanyUser(r)
	rr := record(t, r)
	testRequestStatus(t, rr, 200, "Could not get challenge by key")
}

func TestGetChallengeByKeyUnauthorized(t *testing.T) {
	key := challengeKey.Encode()
	r := recordRequest(t, "GET", "/challenges/"+key, nil)
	testRequestStatus(t, r, 401, "Not logged in user should not be allowed query challenge")
}

func TestGetChallengeByKeyBadKey(t *testing.T) {
	key := "invalidKey"
	r, err := http.NewRequest("GET", "/challenges/"+key, nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	loginAsCompanyUser(r)
	rr := record(t, r)
	testRequestStatus(t, rr, 400, "That key should definitely not have been accepted")
}
