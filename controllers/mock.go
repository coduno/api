package controllers

import (
	"net/http"
	"net/mail"
	"time"

	"github.com/coduno/api/model"
	"github.com/coduno/api/test"
	"github.com/coduno/api/util/password"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func MockChallenge(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	q := model.NewQueryForCompany().Filter("Name =", "Coduno").Limit(1).KeysOnly()

	var companies []model.Company

	keys, err := q.GetAll(ctx, companies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	coduno := keys[0]

	taskOne, _ := model.Task{
		Assignment: model.Assignment{
			Name:         "Hello, world!",
			Description:  "This is the easiest program. It is the hello world of this challenge.",
			Instructions: "Implement a program that outputs \"Hello, world!\" in a programming language of your choice.",
			Duration:     time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "simple-code-task",
			},
		},
		SkillWeights: model.SkillWeights{1, 0, 0},
	}.Put(ctx, nil)

	model.Test{
		Tester: int(test.Simple),
		Name:   "Useless name",
	}.PutWithParent(ctx, taskOne)

	taskTwo, _ := model.Task{
		Assignment: model.Assignment{
			Name:         "Sorting",
			Description:  "This program will require some knowledge about algorithms.",
			Instructions: "Implement a simple bubble sorter on numbers in a programming language of your choice.",
			Duration:     time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "simple-code-task",
			},
		},
		SkillWeights: model.SkillWeights{1, 2, 3},
	}.Put(ctx, nil)

	// Runner:    "simple",
	// Languages: []string{"java", "py"},

	taskThree, _ := model.Task{
		Assignment: model.Assignment{
			Name:         "Some task",
			Description:  "Description of some task",
			Instructions: "Instructions of some task",
			Duration:     time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "simple-code-task",
			},
		},
		SkillWeights: model.SkillWeights{1, 2, 3},
	}.Put(ctx, nil)

	// Runner:    "simple",
	// Languages: []string{"java", "py"},

	model.Challenge{
		Assignment: model.Assignment{
			Name:         "Sequential test",
			Description:  "Description of sequential challenge",
			Instructions: "Instructions of sequential challenge",
			Duration:     time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "sequential-challenge",
			},
		},
		Tasks: []*datastore.Key{taskOne, taskTwo, taskThree},
	}.PutWithParent(ctx, coduno)
}

func Mock(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	pw, _ := password.Hash([]byte("passwordpassword"))

	coduno, _ := model.Company{
		Address: mail.Address{
			Name:    "Coduno",
			Address: "team@cod.uno",
		},
	}.Put(ctx, nil)

	victor, _ := model.User{
		Address: mail.Address{
			Name:    "Victor Balan",
			Address: "victor.balan@cod.uno",
		},
		Nick:           "vbalan",
		HashedPassword: pw,
		Company:        coduno,
	}.Put(ctx, nil)

	paul, _ := model.User{
		Address: mail.Address{
			Name:    "Paul Bochis",
			Address: "paul.bochis@cod.uno",
		},
		Nick:           "pbochis",
		HashedPassword: pw,
		Company:        coduno,
	}.Put(ctx, nil)

	model.User{
		Address: mail.Address{
			Name:    "Alin Mayer",
			Address: "alin.mayer@gmail.com",
		},
		Nick:           "amayer",
		HashedPassword: pw,
	}.Put(ctx, nil)

	lorenz, _ := model.User{
		Address: mail.Address{
			Name:    "Lorenz Leutgeb",
			Address: "lorenz.leutgeb@cod.uno",
		},
		Nick:           "flowlo",
		HashedPassword: pw,
		Company:        coduno,
	}.Put(ctx, nil)

	model.Profile{
		Skills:     model.Skills{12, 40, 1231},
		LastUpdate: time.Now(),
	}.PutWithParent(ctx, victor)

	model.Profile{
		Skills:     model.Skills{11, 1234, 14},
		LastUpdate: time.Now(),
	}.PutWithParent(ctx, paul)

	model.Profile{
		Skills:     model.Skills{154, 12, 1123},
		LastUpdate: time.Now(),
	}.PutWithParent(ctx, lorenz)

	model.User{
		Address: mail.Address{
			Name:    "Admin",
			Address: "admin@cod.uno",
		},
		Nick:           "admin",
		HashedPassword: pw,
	}.PutWithParent(ctx, coduno)

	taskOne, _ := model.Task{
		Assignment: model.Assignment{
			Name:         "Task one",
			Description:  "Description of task one",
			Instructions: "Instructions of task one",
			Duration:     time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "coding-task",
			},
		},
		SkillWeights: model.SkillWeights{1, 2, 3},
	}.Put(ctx, nil)

	// Runner:    "simple",
	// Languages: []string{"java", "py"},

	taskTwo, _ := model.Task{
		Assignment: model.Assignment{
			Name:         "Task two",
			Description:  "Description of task two",
			Instructions: "Instructions of task two",
			Duration:     time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "simple-code-task",
			},
		},
		SkillWeights: model.SkillWeights{1, 2, 3},
	}.Put(ctx, nil)

	// Runner:    "simple",
	// Languages: []string{"java", "py"},

	taskThree, _ := model.Task{
		Assignment: model.Assignment{
			Name:         "Task three",
			Description:  "Description of task three",
			Instructions: "Instructions of task three",
			Duration:     time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "simple-code-task",
			},
		},
		SkillWeights: model.SkillWeights{1, 2, 3},
	}.Put(ctx, nil)

	// Runner:    "simple",
	// Languages: []string{"java", "py"},

	model.Challenge{
		Assignment: model.Assignment{
			Name:         "Challenge one",
			Description:  "Description of challenge one",
			Instructions: "Instructions of challenge one",
			Duration:     time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "sequential-challenge",
			},
		},
		Tasks: []*datastore.Key{taskThree},
	}.PutWithParent(ctx, coduno)

	model.Challenge{
		Assignment: model.Assignment{
			Name:         "Challenge two",
			Description:  "Description of challenge two",
			Instructions: "Instructions of challenge two",
			Duration:     time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "paralel-challenge",
			},
		},
		Tasks: []*datastore.Key{
			taskOne,
			taskTwo,
			taskThree,
		},
	}.PutWithParent(ctx, coduno)
}

func MockCoduno(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	q := model.NewQueryForCompany().Filter("Name =", "Coduno").Limit(1).KeysOnly()

	var companies []model.Company

	keys, err := q.GetAll(ctx, companies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	coduno := keys[0]

	model.Task{
		Assignment: model.Assignment{
			Name:         "Hello, world!",
			Description:  "This is the easiest program. It is the hello world of this challenge.",
			Instructions: "Implement a program that outputs \"Hello, world!\" in a programming language of your choice.",
			Duration:     time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "simple-code-task",
			},
		},
		SkillWeights: model.SkillWeights{1, 0, 0},
	}.PutWithParent(ctx, coduno)

	// Runner:    "simple",
	// Languages: []string{"java", "py", "c", "cpp"},
}
