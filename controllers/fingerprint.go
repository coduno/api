package controllers

import (
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"net/http"

	mailUtils "github.com/coduno/app/mail"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/app/models"
	"github.com/coduno/app/util"
)

// FingerprintData is used to map data from the client.
type FingerprintData struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	ChallengeID string `json:"challangeId"`
	CompanyID   string `json:"companyId"`
}

type MailResponse struct {
	FirstName string
	LastName  string
	Email     string
	Company   string
	Token     string
}

func randomToken() (token string, err error) {
	b := make([]byte, 16)
	_, err = rand.Read(b)
	token = string(b)
	return
}

func HandleFingerprints(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if len(query) == 0 {
		create(w, r, ctx)
		return
	}

	if query["company"][0] == "" {
		http.Error(w, "missing parameter 'company'", http.StatusBadRequest)
		return
	}

	byCompany(query["company"][0], w, r, ctx)
}

func create(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var err error

	if !util.CheckMethod(w, r, "POST") {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var fingerprintData FingerprintData
	if err = json.Unmarshal(body, &fingerprintData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	coder := models.Coder{
		FirstName: fingerprintData.FirstName,
		LastName:  fingerprintData.LastName,
		Email:     fingerprintData.Email,
	}

	coderKey, err := coder.Save(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := randomToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	challengeKey, err := datastore.DecodeKey(fingerprintData.ChallengeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fingerprint := models.Fingerprint{
		Coder:     coderKey,
		Challenge: challengeKey,
		Token:     token,
	}

	companyKey, err := datastore.DecodeKey(fingerprintData.CompanyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// TODO(pbochis): This is where we will send an e-mail to the candidate with
	// somthing like "cod.uno/fingerprint/:token".

	var company models.Company
	err = datastore.Get(ctx, companyKey, &company)

	mailResponse := MailResponse{
		FirstName: fingerprintData.FirstName,
		LastName:  fingerprintData.LastName,
		Email:     fingerprintData.Email,
		Company:   company.Name,
		Token:     token}
	err = sendMailToCoder(ctx, mailResponse)
	if err != nil {
		http.Error(w, "Could not send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	key, err := fingerprint.Save(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.WriteEntity(w, key, fingerprint)
}

func byCompany(companyKey string, w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if !util.CheckMethod(w, r, "GET") {
		return
	}

	key, _ := datastore.DecodeKey(companyKey)
	q := datastore.NewQuery(models.ChallengeKind).Filter("Company = ", key).KeysOnly()
	keys, err := q.GetAll(ctx, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q = datastore.NewQuery(models.FingerprintKind)
	for _, val := range keys {
		q.Filter("Challenge = ", val)
	}
	var fingerprints []models.Fingerprint
	keys, err = q.GetAll(ctx, &fingerprints)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	values := make([]interface{}, len(fingerprints))
	for i, fingerprint := range fingerprints {
		values[i] = fingerprint
	}
	util.WriteEntities(w, keys, values)
}

func sendMailToCoder(c context.Context, mailResponse MailResponse) error {
	body, err := mailUtils.PrepareMailTemplate(mailUtils.CandidateInvitationTemplate, mailResponse)
	if err != nil {
		return err
	}
	return mailUtils.SendMail(c, []string{mailResponse.Email}, "Your coding challenge", body)
}
