package controllers

import (
	"net/http"
	"net/mail"
	"time"

	"github.com/coduno/app/model"
	"github.com/coduno/engine/util/password"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func Mock(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	pw, _ := password.Hash([]byte("password"))

	coduno, _ := model.Company{
		Address: mail.Address{
			Name:    "Coduno",
			Address: "team@cod.uno",
		},
	}.Save(ctx)

	victor, _ := model.User{
		Address: mail.Address{
			Name:    "Victor Balan",
			Address: "victor.balan@cod.uno",
		},
		Nick:           "vbalan",
		HashedPassword: pw,
	}.SaveWithParent(ctx, coduno)

	paul, _ := model.User{
		Address: mail.Address{
			Name:    "Paul Bochis",
			Address: "paul.bochis@cod.uno",
		},
		Nick:           "pbhcis",
		HashedPassword: pw,
	}.SaveWithParent(ctx, coduno)

	lorenz, _ := model.User{
		Address: mail.Address{
			Name:    "Lorenz Leutgeb",
			Address: "lorenz.leutgeb@cod.uno",
		},
		Nick:           "flowlo",
		HashedPassword: pw,
	}.SaveWithParent(ctx, coduno)

	model.Profile{
		Skills:     model.Skills{12, 40, 1231},
		LastUpdate: time.Now(),
	}.SaveWithParent(ctx, victor)

	model.Profile{
		Skills:     model.Skills{11, 1234, 14},
		LastUpdate: time.Now(),
	}.SaveWithParent(ctx, paul)

	model.Profile{
		Skills:     model.Skills{154, 12, 1123},
		LastUpdate: time.Now(),
	}.SaveWithParent(ctx, lorenz)

	model.User{
		Address: mail.Address{
			Name:    "Admin",
			Address: "admin@cod.uno",
		},
		Nick:           "admin",
		HashedPassword: pw,
	}.SaveWithParent(ctx, coduno)

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
	}.Save(ctx)

	taskTwo, _ := model.Task{
		Assignment: model.Assignment{
			Name:         "Task two",
			Description:  "Description of task two",
			Instructions: "Instructions of task two",
			Duration:     time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "boolean-task",
			},
		},
		SkillWeights: model.SkillWeights{1, 2, 3},
	}.Save(ctx)

	taskThree, _ := model.CodeTask{
		Task: model.Task{
			Assignment: model.Assignment{
				Name:         "Task three",
				Description:  "Description of task three",
				Instructions: "Instructions of task three",
				Duration:     time.Hour,
				Endpoints: model.Endpoints{
					WebInterface: "multiple-select-task",
				},
			},
			SkillWeights: model.SkillWeights{1, 2, 3},
		},
		Runner: "simple",
	}.Save(ctx)

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
		Tasks: []*datastore.Key{taskOne},
	}.SaveWithParent(ctx, coduno)

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
	}.SaveWithParent(ctx, coduno)
}
