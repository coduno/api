package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util/password"
	"github.com/coduno/engine/passenger"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

var instance aetest.Instance

func record(t *testing.T, r *http.Request) *httptest.ResponseRecorder {
	instance.Associate(r)
	w := httptest.NewRecorder()
	Handler().ServeHTTP(w, r)
	return w
}

func testRequestStatus(t *testing.T, r *httptest.ResponseRecorder, expected int, message string) {
	if expected == r.Code {
		return
	}
	t.Fatalf("Expected HTTP status %d but got %d: %s", expected, r.Code, message)
}

func backgroundContext() context.Context {
	r, err := instance.NewRequest("GET", "/", strings.NewReader(""))
	if err != nil {
		panic(err)
	}
	return appengine.NewContext(r)
}

func recordRequest(t *testing.T, method, route string, body interface{}) *httptest.ResponseRecorder {
	r, err := http.NewRequest(method, route, requestBody(t, body))
	if err != nil {
		t.Fatal(err)
	}
	return record(t, r)
}

var CompanyKey *datastore.Key
var CompanyUserKey *datastore.Key
var UserKey *datastore.Key
var TaskKey *datastore.Key
var challengeKey *datastore.Key

const companyUserName = "john"
const coderUserName = "andy"
const testingPassword = "passwordpassword"

func mockData(ctx context.Context) {
	var err error
	CompanyKey, err = model.Company{
		Address: mail.Address{
			Name:    "Coduno",
			Address: "team@cod.uno",
		},
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	pw, _ := password.Hash([]byte(testingPassword))
	CompanyUserKey, err = model.User{
		Address: mail.Address{
			Name:    companyUserName,
			Address: "john@cod.uno",
		},
		Nick:           "john",
		HashedPassword: pw,
		Company:        CompanyKey,
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	UserKey, err = model.User{
		Address: mail.Address{
			Name:    coderUserName,
			Address: "andy@example.com",
		},
		Nick:           "andy",
		HashedPassword: pw,
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	TaskKey, err = model.Task{
		Assignment: model.Assignment{
			Name:         "Task one",
			Description:  "Description of task one",
			Instructions: "Instructions of task one",
			Duration:     time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "coding-task",
			},
		},
		Languages:    []string{"py", "java"},
		SkillWeights: model.SkillWeights{0.25, 0.25, 0.25, 0.25},
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	challengeKey, err = model.Challenge{
		Assignment: model.Assignment{
			Name:         "Challenge one",
			Description:  "Description of challenge one",
			Instructions: "Instructions of challenge one yay",
			Duration:     time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "sequential-challenge",
			},
		},
		Tasks: []*datastore.Key{TaskKey},
	}.PutWithParent(ctx, CompanyKey)
	if err != nil {
		panic(err)
	}

	var cmp model.User
	datastore.Get(ctx, CompanyUserKey, &cmp)
}

func loginAsCompanyUser(r *http.Request) {
	r.SetBasicAuth(companyUserName, testingPassword)
}

func LoginAsCoderUser(t *testing.T, ctx context.Context, r *http.Request) context.Context {
	r.SetBasicAuth(coderUserName, testingPassword)
	c, err := passenger.NewContextFromRequest(ctx, r)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// CreateAndSaveCompany creates a new Company and saves it to datastore using the provided context
func CreateAndSaveCompany(t *testing.T, ctx context.Context, name, email string) (model.Company, *datastore.Key) {
	if name == "" {
		name = "Example Company"
	}
	if email == "" {
		email = "office@example.com"
	}
	company := model.Company{
		Address: mail.Address{
			Name:    name,
			Address: email,
		},
	}
	var err error
	var key *datastore.Key

	key, err = company.Put(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	return company, key
}

func CreateAndSaveChallenge(t *testing.T, ctx context.Context, company *datastore.Key) (model.Challenge, *datastore.Key) {
	assignment := model.Assignment{
		Name:         "Challenge Name",
		Description:  "Challenge Description",
		Instructions: "Challenge Instructions",
		Duration:     time.Hour,
		Endpoints: model.Endpoints{
			WebInterface: "sequential-challenge",
		},
	}
	tasks := []*datastore.Key{TaskKey}
	challenge := model.Challenge{
		Assignment: assignment,
		Tasks:      tasks,
	}
	var err error
	var key *datastore.Key
	if key, err = challenge.PutWithParent(ctx, company); err != nil {
		t.Fatal(err)
	}
	return challenge, key
}

func requestBody(t *testing.T, data interface{}) io.Reader {
	if data == nil {
		return nil
	}
	var jsonData, err = json.Marshal(data)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return bytes.NewReader(jsonData)
}

func TestMain(m *testing.M) {
	// FIXME: Tests shouldn't depend on templates belonging to the app in production.
	InvitationTemplatePath = "../mail/template.invitation"
	SubTemplatePath = "../mail/template.subscription"

	var err error
	instance, err = aetest.NewInstance(&aetest.Options{
		AppID: "coduno",
		StronglyConsistentDatastore: true,
	})
	if err != nil {
		panic(err)
	}
	mockData(backgroundContext())
	exitCode := m.Run()

	// Tear down testing instance.
	instance.Close()

	os.Exit(exitCode)
}
