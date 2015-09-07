package controllers

import (
	"testing"

	"github.com/coduno/api/tests/testUtils"
)

func TestGetChallengeByKey(t *testing.T) {
	key := testUtils.ChallengeKey.Encode()
	req, w := requestAndResponse(t, "GET", "/challenges/"+key, nil)
	ctx = testUtils.LoginAsCompanyUser(t, ctx, req)
	r.ServeHTTP(w, req)
	testRequestStatus(t, 200, "Could not get challenge by key")
	ctx = testUtils.Logout(ctx)
}

func TestGetChallengeByKeyUnauthorized(t *testing.T) {
	key := testUtils.ChallengeKey.Encode()
	req, w := requestAndResponse(t, "GET", "/challenges/"+key, nil)
	r.ServeHTTP(w, req)
	testRequestStatus(t, 401, "Not logged in user should not be allowed query challenge")
}

func TestGetChallengeByKeyBadKey(t *testing.T) {
	key := "thisShouldDeffinetlyNotWork"
	req, w := requestAndResponse(t, "GET", "/challenges/"+key, nil)
	ctx = testUtils.LoginAsCompanyUser(t, ctx, req)
	r.ServeHTTP(w, req)
	testRequestStatus(t, 400, "That key should deffinetly not have been accepted")
}
