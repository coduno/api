package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/coduno/api/tests/testUtils"
	"github.com/coduno/api/util/routing"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
)

var ctx context.Context
var instance aetest.Instance
var r *mux.Router
var lastStatus int
var lastError error

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/companies", testSetup(PostCompany))
	r.HandleFunc("/challenges/{key}", testSetup(ChallengeByKey))
	http.Handle("/", r)
	return r
}

func testSetup(h routing.ContextHandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastStatus, lastError = h(ctx, w, r)
	})
}

func testRequestStatusAndError(t *testing.T, expectedStatus int, expectedError error) {
	if expectedStatus != lastStatus {
		t.Fatal("Unexpected return status")
	}
	if expectedError != lastError {
		t.Fatal(lastError)
	}
	lastError = nil
}

func testRequestStatus(t *testing.T, expectedStatus int, message string) {
	if expectedStatus != lastStatus {
		t.Fatal(message)
	}
	lastError = nil
}

func requestAndResponse(t *testing.T, method string, route string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {
	req, err := instance.NewRequest(method, route, testUtils.RequestBody(t, body))
	if err != nil {
		t.Fatal(err)
	}
	return req, httptest.NewRecorder()
}

func TestMain(m *testing.M) {
	// FIXME: Tests shouldn't depend on templates belonging to the app in production.
	InvitationTemplatePath = "../mail/template.invitation"
	SubTemplatePath = "../mail/template.subscription"

	var err error
	var callback func()
	ctx, callback, err = aetest.NewContext()
	if err != nil {
		os.Exit(1)
	}
	defer callback()

	instance, err = aetest.NewInstance(nil)
	if err != nil {
		os.Exit(1)
	}
	testUtils.MockData(ctx)
	r = Router()
	os.Exit(m.Run())
}
