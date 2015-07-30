package controllers

import (
	"net/http"
	"net/mail"
	"time"

	. "github.com/coduno/engine/model"
	"github.com/coduno/engine/util/password"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func MockData(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	pw, _ := password.Hash([]byte("password"))
	victor, _ := User{Address: mail.Address{Name: "Victor Balan", Address: "victor.balan@cod.uno"}, Nick: "vbalan", HashedPassword: pw}.Save(ctx)
	paul, _ := User{Address: mail.Address{Name: "Paul Bochis", Address: "paul.bochis@cod.uno"}, Nick: "pbhcis", HashedPassword: pw}.Save(ctx)
	lorenz, _ := User{Address: mail.Address{Name: "Lorenz Leutgeb", Address: "lorenz.leutgeb@cod.uno"}, Nick: "flowlo", HashedPassword: pw}.Save(ctx)

	Profile{Skills: Skills{12, 40, 1231}, LastUpdate: time.Now()}.SaveWithParent(ctx, victor)
	Profile{Skills: Skills{11, 1234, 14}, LastUpdate: time.Now()}.SaveWithParent(ctx, paul)
	Profile{Skills: Skills{154, 12, 1123}, LastUpdate: time.Now()}.SaveWithParent(ctx, lorenz)

	coduno, _ := Company{Address: mail.Address{Name: "Coduno", Address: "office@cod.uno"}}.Save(ctx)
	User{Address: mail.Address{Name: "Admin", Address: "admin@cod.uno"}, Nick: "codunoadmin", HashedPassword: pw}.SaveWithParent(ctx, coduno)

	taskOne, _ := Task{Assignment: Assignment{Name: "Task one", Description: "Description of task one", Instructions: "Instructions of task one",
		Duration: time.Hour, Endpoints: Endpoints{WebInterface: "coding-task"}},
		SkillWeights: SkillWeights{1, 2, 3}}.Save(ctx)

	taskTwo, _ := Task{Assignment: Assignment{Name: "Task two", Description: "Description of task two", Instructions: "Instructions of task two",
		Duration: time.Hour, Endpoints: Endpoints{WebInterface: "boolean-task"}},
		SkillWeights: SkillWeights{1, 2, 3}}.Save(ctx)

	taskThree, _ := Task{Assignment: Assignment{Name: "Task three", Description: "Description of task three", Instructions: "Instructions of task three",
		Duration: time.Hour, Endpoints: Endpoints{WebInterface: "multiple-select-task"}},
		SkillWeights: SkillWeights{1, 2, 3}}.Save(ctx)

	var cOneTasks = make([]*datastore.Key, 1)
	cOneTasks[0] = taskOne

	Challenge{Assignment: Assignment{Name: "Challenge one", Description: "Description of challenge one",
		Instructions: "Instructions of challenge one", Duration: time.Hour,
		Endpoints: Endpoints{WebInterface: "secvential-challenge"}}, Tasks: cOneTasks}.SaveWithParent(ctx, coduno)

	var cTwoTasks = make([]*datastore.Key, 3)
	cTwoTasks[0] = taskOne
	cTwoTasks[1] = taskTwo
	cTwoTasks[2] = taskThree
	Challenge{Assignment: Assignment{Name: "Challenge two", Description: "Description of challenge two",
		Instructions: "Instructions of challenge two", Duration: time.Hour,
		Endpoints: Endpoints{WebInterface: "paralel-challenge"}}, Tasks: cTwoTasks}.SaveWithParent(ctx, coduno)

}
