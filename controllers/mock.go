package controllers

import (
	"net/http"
	"net/mail"
	"time"

	"golang.org/x/net/context"

	"github.com/coduno/api/logic"
	"github.com/coduno/api/model"
	"github.com/coduno/api/test"
	"github.com/coduno/api/util"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func init() {
	router.HandleFunc("/mock", Mock)
}

func Mock(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

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
		HashedPassword: []byte{0x24, 0x32, 0x61, 0x24, 0x31, 0x30, 0x24, 0x42, 0x2e, 0x79, 0x5a, 0x2f, 0x4f, 0x6e, 0x41, 0x4d, 0x47, 0x71, 0x6f, 0x51, 0x76, 0x41, 0x61, 0x39, 0x49, 0x53, 0x79, 0x38, 0x2e, 0x5a, 0x4d, 0x2e, 0x38, 0x6d, 0x31, 0x41, 0x70, 0x4a, 0x45, 0x46, 0x48, 0x4c, 0x70, 0x5a, 0x75, 0x59, 0x6f, 0x56, 0x48, 0x67, 0x6e, 0x63, 0x34, 0x50, 0x6b, 0x42, 0x70, 0x47, 0x78, 0x4b},
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
		HashedPassword: []byte{0x24, 0x32, 0x61, 0x24, 0x31, 0x30, 0x24, 0x5a, 0x6c, 0x6f, 0x4e, 0x57, 0x46, 0x6d, 0x6a, 0x6a, 0x73, 0x76, 0x71, 0x35, 0x55, 0x6b, 0x44, 0x36, 0x4f, 0x6e, 0x74, 0x49, 0x2e, 0x47, 0x75, 0x47, 0x49, 0x33, 0x6f, 0x6e, 0x43, 0x53, 0x59, 0x53, 0x56, 0x6c, 0x36, 0x6e, 0x59, 0x50, 0x70, 0x4c, 0x55, 0x71, 0x61, 0x6e, 0x53, 0x77, 0x37, 0x70, 0x64, 0x4b, 0x37, 0x53},
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
		HashedPassword: []byte{0x24, 0x32, 0x61, 0x24, 0x31, 0x30, 0x24, 0x78, 0x4a, 0x2f, 0x4a, 0x65, 0x57, 0x74, 0x46, 0x33, 0x55, 0x72, 0x2e, 0x36, 0x59, 0x75, 0x35, 0x6f, 0x38, 0x52, 0x77, 0x47, 0x75, 0x32, 0x4a, 0x35, 0x47, 0x69, 0x58, 0x67, 0x55, 0x4b, 0x72, 0x68, 0x51, 0x4d, 0x4d, 0x61, 0x72, 0x75, 0x47, 0x65, 0x36, 0x2e, 0x69, 0x34, 0x73, 0x39, 0x73, 0x7a, 0x54, 0x70, 0x63, 0x79},
		Company:        coduno,
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	model.Profile{
		Skills:     model.Skills{},
		LastUpdate: time.Now(),
	}.PutWithParent(ctx, victor)

	model.Profile{
		Skills:     model.Skills{},
		LastUpdate: time.Now(),
	}.PutWithParent(ctx, paul)

	model.Profile{
		Skills:     model.Skills{},
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
		Languages: []string{"java", "py", "c", "cpp"},
		SkillWeights: model.SkillWeights{
			Algorithmics: 0.1,
			Readability:  0.1,
			Security:     0.1,
			CodingSpeed:  0.7,
		},
		Tasker: int64(logic.DiffTasker),
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	_, err = model.Test{
		Tester: int64(test.Diff),
		Name:   "Hello world test",
		Params: map[string]string{
			// TODO(victorbalan): Extract params in constants
			"tests": "helloworld/helloworld",
		},
	}.PutWithParent(ctx, taskOne)
	if err != nil {
		panic(err)
	}

	taskTwo, err := model.Task{
		Assignment: model.Assignment{
			Name: "Fizz Buzz",
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
		SkillWeights: model.SkillWeights{
			Algorithmics: 0.1,
			Readability:  0.2,
			Security:     0,
			CodingSpeed:  0.7,
		},
		Languages: []string{"java", "py", "c", "cpp"},
		Tasker:    int64(logic.DiffTasker),
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	model.Test{
		Tester: int64(test.IO),
		Params: map[string]string{
			"bucket": "coduno-tests",
			"input":  "fizzbuzz/fizzbuzzin10^2",
			"output": "fizzbuzz/fizzbuzz10^2",
		},
	}.PutWithParent(ctx, taskTwo)

	model.Test{
		Tester: int64(test.IO),
		Params: map[string]string{
			"bucket": "coduno-tests",
			"input":  "fizzbuzz/fizzbuzzin10^3",
			"output": "fizzbuzz/fizzbuzz10^3",
		},
	}.PutWithParent(ctx, taskTwo)

	model.Test{
		Tester: int64(test.IO),
		Params: map[string]string{
			"bucket": "coduno-tests",
			"input":  "fizzbuzz/fizzbuzzin10^4",
			"output": "fizzbuzz/fizzbuzz10^4",
		},
	}.PutWithParent(ctx, taskTwo)

	taskThree, err := model.Task{
		Assignment: model.Assignment{
			Name: "N-Gram",
			Description: `In the fields of computational linguistics and probability, an n-gram is a contiguous sequence
			of n items from a given sequence of text or speech. The items can be phonemes, syllables, letters, words or base
			pairs according to the application. The n-grams typically are collected from a text or speech corpus.`,
			Instructions: "Your job is to create a function with the signature `int ngram(String text, int len)` and outputs the number of n-grams of length `len`.",
			Duration:     time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "javaut-task",
			},
		},
		SkillWeights: model.SkillWeights{
			Algorithmics: 0.1,
			Readability:  0.1,
			Security:     0.1,
			CodingSpeed:  0.7,
		},
		Languages: []string{"java"},
		Tasker:    int64(logic.JunitTasker),
		Templates: templateHelper(map[string][]string{"java": {"ngram/Application.java"}}),
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	model.Test{
		Tester: int64(test.Junit),
		Params: map[string]string{
			"test":        "ngram/Tests.java",
			"resultPath":  "/run/build/test-results/",
			"imageSuffix": "javaut",
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
		SkillWeights: model.SkillWeights{},
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
		Resulter: int64(logic.Average),
	}.PutWithParent(ctx, coduno)
	if err != nil {
		panic(err)
	}

	taskPrimeUT, err := model.Task{
		Assignment: model.Assignment{
			Name:         "Primes test",
			Description:  `Test a method which checks wether an integer is prime.`,
			Instructions: "You have to write java unit tests in order to check wether the method Application.isPrime(int n) returns the correct answer.",
			Duration:     time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "javaut-task",
			},
		},
		SkillWeights: model.SkillWeights{
			Algorithmics: 0.4,
			Readability:  0.3,
			CodingSpeed:  0.3,
			Security:     0,
		},
		Languages: []string{"java"},
		Tasker:    int64(logic.JunitTasker),
		Templates: templateHelper(map[string][]string{"java": {"primes/Application.java"}}),
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	model.Test{
		Name:   "Correct alg",
		Tester: int64(test.CoderJunit),
		Params: map[string]string{
			"code":        "primes/v1/",
			"resultPath":  "/run/build/test-results/",
			"imageSuffix": "javaut",
			"shouldFail":  "false",
		},
	}.PutWithParent(ctx, taskPrimeUT)

	model.Test{
		Name:   "Broken alg",
		Tester: int64(test.CoderJunit),
		Params: map[string]string{
			"code":        "primes/v2/",
			"resultPath":  "/run/build/test-results/",
			"imageSuffix": "javaut",
			"shouldFail":  "true",
		},
	}.PutWithParent(ctx, taskPrimeUT)

	MockFrequentisChallenge(ctx, coduno, w, req)
}

func MockFrequentisChallenge(ctx context.Context, coduno *datastore.Key, w http.ResponseWriter, req *http.Request) {
	taskOne, err := model.Task{
		Assignment: model.Assignment{
			Name:        "CoinBot",
			Description: "CoinBot is a simple remote-controlled robot. He is placed in a big hall, that contains coins. His task is to collect all coins and return to the position he started from.",
			Instructions: `You are given an overview of the hall that ConBot is placed in.
			  	Furthermore, you'll controll him by issuing simple commands. You can tell CoinBot to move forward, by using the "move n"
			   	command where n is the number of fields he should move. Also, you can make him turn left or right
		 			by using the commands "left" and "right" respectively. If CoinBot is at the location of a coin, you need to
		 			instruct him to pick up the coin with the command "pick". You can make CoinBot carry out your commands by hitting the arrow.
				 	If you do this repeatedly, it will cause a reset of the game, so that you always start with the same environment.
			 		Once CoinBot has fulfilled his mission, you will automatically advance to the next task.`,
			Duration: time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "canvas-game-task",
			},
		},
		SkillWeights: model.SkillWeights{
			Algorithmics: 0.1,
			Readability:  0.1,
			Security:     0.1,
			CodingSpeed:  0.7,
		},
		Templates: templateHelper(map[string][]string{"json": {"robot/robot.json"}}),
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	_, err = model.Test{
		Tester: int64(test.Robot),
		Name:   "Robot test",
		Params: map[string]string{
			"tests": "robot/robot.json",
		},
	}.PutWithParent(ctx, taskOne)
	if err != nil {
		panic(err)
	}

	taskTwo, err := model.Task{
		Assignment: model.Assignment{
			Name:        "AvlTree",
			Description: "AvlTree",
			Instructions: `Your task is to write Junit unit tests for an [AvlTree](https://en.wikipedia.org/wiki/AVL_tree) implementation.
			Below you are given the signatures and descriptions of all public operations of the AvlTree class.
			 * void insert(int k)	 	Insert k if it doesn't exist. Duplicates will be ignored.
			 * void remove(int k)	 	Remove x if it exists.
			 * int size()		 	 	Return number of vertexes.
			 * boolean contains(int k) 	Returns true if a vertex with value k exists in the tree, false otherwise.
			 * int findMinimum()    	Returns the smallest value in the tree. If the tree is empty, will return Integer.MIN_VALUE.
			 * int findMaximum()    	Returns the highest value in the tree. If the tree is empty, will return Integer.MIN_VALUE.
			 * boolean isEmpty( )       Return true if the tree is empty, false otherwise.
			 * void empty( )    		Remove all items from the tree.
			`,
			Duration: time.Hour,
			Endpoints: model.Endpoints{
				WebInterface: "coder-javaut-task",
			},
		},
		SkillWeights: model.SkillWeights{
			Algorithmics: 0.4,
			Readability:  0.2,
			Security:     0,
			CodingSpeed:  0.4,
		},
		Templates: templateHelper(map[string][]string{"java": {"avl/Tests.java"}}),
		Languages: []string{"java"},
		Tasker:    int64(logic.JunitTasker),
	}.Put(ctx, nil)
	if err != nil {
		panic(err)
	}

	testsForTaskTwo(ctx, taskTwo)

	_, err = model.Challenge{
		Assignment: model.Assignment{
			Name:        "Frequentis hiring challenge",
			Description: "This is a hiring challenge for the Frequentis company.",
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
		},
		Resulter: int64(logic.Average),
	}.PutWithParent(ctx, coduno)
	if err != nil {
		panic(err)
	}
}

func testsForTaskTwo(ctx context.Context, taskTwo *datastore.Key) {
	model.Test{
		Name:   "v1",
		Tester: int64(test.CoderJunit),
		Params: map[string]string{
			"code":        "avl/v1/",
			"resultPath":  "/run/build/test-results/",
			"imageSuffix": "javaut",
			"shouldFail":  "false",
		},
	}.PutWithParent(ctx, taskTwo)

	for _, v := range []string{"v2", "v3", "v4", "v5", "v6", "v7"} {
		model.Test{
			Name:   v,
			Tester: int64(test.CoderJunit),
			Params: map[string]string{
				"code":        "avl/" + v + "/",
				"resultPath":  "/run/build/test-results/",
				"imageSuffix": "javaut",
				"shouldFail":  "true",
			},
		}.PutWithParent(ctx, taskTwo)
	}
}

func templateHelper(m map[string][]string) map[string][]model.StoredObject {
	res := map[string][]model.StoredObject{}
	for k, files := range m {
		sos := make([]model.StoredObject, 0, len(files))
		for _, file := range files {
			sos = append(sos, model.StoredObject{
				Bucket: util.TemplateBucket,
				Name:   file,
			})
		}
		res[k] = sos
	}
	return res
}
