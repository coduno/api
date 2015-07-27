package controllers

import (
	"net/http"

	"github.com/coduno/engine/appengine/model"
	"github.com/coduno/engine/util/password"
	"google.golang.org/appengine"
)

func MockCompany(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	pw, _ := password.Hash([]byte("123123123123"))
	cmp := model.Company{Name: "cat", Email: "paul@cod.uno", HashedPassword: pw}
	cmp.Save(ctx)
	w.Write([]byte("Its all fine"))
}

func MockData(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	company := model.Company{Name: "Catalysts"}
	companyKey, _ := company.Save(ctx)

	challenge := model.Challenge{Name: "Tic-Tac-Toe", Instructions: "Implenet tic tac toe input and output blah blah", Company: companyKey}
	challengeKey, _ := challenge.Save(ctx)

	template := model.Template{Language: "Java", Path: "/templates/TicTacToeTemplate.java", Challenge: challengeKey}
	template.Save(ctx)

	coder := model.Coder{Email: "victor.balan@cod.uno", FirstName: "Victor", LastName: "Balan"}
	coderKey, _ := coder.Save(ctx)

	fingerprint := model.Fingerprint{Coder: coderKey, Challenge: challengeKey, Token: "deadbeefcafebabe"}
	fingerprint.Save(ctx)
}
