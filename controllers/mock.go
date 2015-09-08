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

func Mock(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	pw, _ := password.Hash([]byte("passwordpassword"))

	coduno, err := model.Company{
		Address: mail.Address{
			Name:    "Coduno",
			Address: "team@cod.uno",
		},
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	victor, err := model.User{
		Address: mail.Address{
			Name:    "Victor Balan",
			Address: "victor.balan@cod.uno",
		},
		Nick:           "vbalan",
		HashedPassword: pw,
		Company:        coduno,
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	paul, err := model.User{
		Address: mail.Address{
			Name:    "Paul Bochis",
			Address: "paul.bochis@cod.uno",
		},
		Nick:           "pbochis",
		HashedPassword: pw,
		Company:        coduno,
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	lorenz, err := model.User{
		Address: mail.Address{
			Name:    "Lorenz Leutgeb",
			Address: "lorenz.leutgeb@cod.uno",
		},
		Nick:           "flowlo",
		HashedPassword: pw,
		Company:        coduno,
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	model.Profile{
		Skills:     model.Skills{1, .5, 1},
		LastUpdate: time.Now(),
	}.PutWithParent(ctx, victor)

	model.Profile{
		Skills:     model.Skills{.5, 1, 1},
		LastUpdate: time.Now(),
	}.PutWithParent(ctx, paul)

	model.Profile{
		Skills:     model.Skills{1, 1, .5},
		LastUpdate: time.Now(),
	}.PutWithParent(ctx, lorenz)

	taskOne, err := model.Task{
		Assignment: model.Assignment{
			Name:         "Hello, world!",
			Description:  "This is a welcome task to our platform. It is the easiest one so you can learn the ui and the workflow.",
			Instructions: "Create a program that outputs 'Hello, world!' in a language of your preference.",
			Duration:     time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "output-match-task",
			},
		},
		Languages:    []string{"java", "py", "c", "cpp"},
		SkillWeights: model.SkillWeights{.1, .2, .3},
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	_, err = model.Test{
		Tester: int64(test.Simple),
		Name:   "Hello world test",
		Params: map[string]string{
			"tests": "coduno-tests/helloworld",
		},
	}.PutWithParent(ctx, taskOne)
	if err != nil {
		panic(err)
	}

	taskTwo, err := model.Task{
		Assignment: model.Assignment{
			Name: "Fizzbuzz",
			Description: `Fizz buzz is a group word game for children to teach them about division.
			 Players take turns to count incrementally, replacing any number divisible by three with the word 'fizz',
			 and any number divisible by five with the word 'buzz'.`,
			Instructions: `Your job is to create the 'fizzbuzz(int n)' function.
			The n parameter represents the max number to wich you need to generate the fizzbuzz data.
			The output needs to be separated by '\n'.`,
			Duration: time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "output-match-task",
			},
		},
		SkillWeights: model.SkillWeights{.1, .2, .3},
		Languages:    []string{"java", "py", "c", "cpp"},
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	model.Test{
		Tester: int64(test.Simple),
		Params: map[string]string{
			"tests": "coduno-tests/fizzbuzz",
		},
	}.PutWithParent(ctx, taskTwo)

	taskThree, err := model.Task{
		Assignment: model.Assignment{
			Name: "N-Gram",
			Description: `In the fields of computational linguistics and probability, an n-gram is a contiguous sequence
			of n items from a given sequence of text or speech. The items can be phonemes, syllables, letters, words or base
			pairs according to the application. The n-grams typically are collected from a text or speech corpus.`,
			Instructions: `Your job is to create a function with the signature ngram(String content, int len)
			and outputs the number of ngrams of length len.`,
			Duration: time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "javaut-task",
			},
		},
		SkillWeights: model.SkillWeights{.1, .2, .3},
		Languages:    []string{"javaut"},
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	model.Test{
		Tester: int64(test.Junit),
		Params: map[string]string{
			"tests":      "coduno-tests/ngram",
			"resultPath": "/run/build/test-results/",
		},
	}.PutWithParent(ctx, taskThree)

	taskFour, err := model.Task{
		Assignment: model.Assignment{
			Name:         "Simple code run test",
			Description:  "This is a mocked task for testing the simple code run.",
			Instructions: "This task will not be tested against anything.",
			Duration:     time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "simple-code-task",
			},
		},
		SkillWeights: model.SkillWeights{.1, .2, .3},
		Languages:    []string{"java", "py", "c", "cpp"},
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	model.Test{
		Tester: int64(test.Simple),
	}.PutWithParent(ctx, taskFour)

	_, err = model.Challenge{
		Assignment: model.Assignment{
			Name:        "Coduno hiring challenge",
			Description: "This is a hiring challenge for the Coduno company.",
			Instructions: `You can select your preffered language from the languages
			dropdown at every run your code will be tested so be careful with what you run.
			You can finish anytime and start the next task but keep in mind that you will not be
			able to get back to the previous task. Good luck!`,
			Duration: time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "sequential-challenge",
			},
		},
		Tasks: []*datastore.Key{
			taskOne,
			taskTwo,
			taskThree,
		},
	}.PutWithParent(ctx, coduno)
	if err != nil {
		panic(err)
	}
}
