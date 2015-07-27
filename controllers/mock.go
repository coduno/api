package controllers

import (
	"net/http"

	"github.com/coduno/app/models"
	"github.com/coduno/app/util/password"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func MockCompany(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	pw, _ := password.Hash([]byte("123123123123"))
	cmp := models.Company{Name: "cat", Email: "paul@cod.uno", HashedPassword: pw}
	cmp.Save(ctx)
	w.Write([]byte("Its all fine"))
}

func MockData(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	company := models.Company{Name: "Catalysts"}
	companyKey, _ := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "companies", nil), &company)

	challenge := models.Challenge{Name: "Tic-Tac-Toe", Instructions: "Implenet tic tac toe input and output blah blah", Company: companyKey}
	challengeKey, _ := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "challenges", nil), &challenge)

	template := models.Template{Language: "Java", Path: "/templates/TicTacToeTemplate.java", Challenge: challengeKey}
	datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "templates", nil), &template)

	coder := models.Coder{Email: "victor.balan@cod.uno", FirstName: "Victor", LastName: "Balan"}
	coderKey, _ := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "coders", nil), &coder)

	fingerprint := models.Fingerprint{Coder: coderKey, Challenge: challengeKey, Token: "deadbeefcafebabe"}
	datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "fingerprints", nil), &fingerprint)
}
