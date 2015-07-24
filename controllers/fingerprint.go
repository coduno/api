package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/app/models"
	"github.com/coduno/app/util"
	"github.com/gorilla/mux"
	"github.com/m4rw3r/uuid"
)

//FingerprintData  is used to map data from the client.
type FingerprintData struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	ChallangeID string `json:"challangeId"`
}

//CreateFingerprint from the request body
func CreateFingerprint(w http.ResponseWriter, r *http.Request, ctx context.Context) {

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
		Email:     fingerprintData.Email}

	coderKey, err := coder.Save(ctx)

	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := uuid.V4()

	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	challangeKey, err := datastore.DecodeKey(fingerprintData.ChallangeID)

	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fingerprint := models.Fingerprint{Coder: coderKey, Challenge: challangeKey, Token: token.String()}

	_, err = fingerprint.Save(ctx)

	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
	}

	response, err := json.Marshal(fingerprint)

	if err != nil {
		http.Error(w, "Json marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO(pbochis): This is where we will send an e-mail to the candidate with
	// somthing like "cod.uno/fingerprint/:token".

	w.Write([]byte(response))
	return
}

// LoadFingerprintsByCompanyID -
func LoadFingerprintsByCompanyID(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if !util.CheckMethod(w, r, "GET") {
		return
	}

	companyID := mux.Vars(r)["companyId"]
	key, _ := datastore.DecodeKey(companyID)
	q := datastore.NewQuery("challenges").Filter("Company = ", key).KeysOnly()
	keys, err := q.GetAll(ctx, nil)
	if err != nil {
		http.Error(w, "Datastore error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	q = datastore.NewQuery("fingerprints")
	for _, val := range keys {
		q.Filter("Challenge = ", val)
	}
	var fingerprints []models.Fingerprint
	keys, err = q.GetAll(ctx, &fingerprints)
	if err != nil {
		http.Error(w, "Datastore error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	for i := 0; i < len(keys); i++ {
		fingerprints[i].EntityID = keys[i].Encode()
	}
	json, err := json.Marshal(fingerprints)
	if err != nil {
		http.Error(w, "Json marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(json))
}
