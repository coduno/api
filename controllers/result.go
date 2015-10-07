package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/coduno/api/logic"
	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"github.com/coduno/api/util/passenger"
	"github.com/gorilla/mux"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"golang.org/x/net/context"
)

func init() {
	router.Handle("/results", ContextHandlerFunc(CreateResult))
	router.Handle("/results/{resultKey}", ContextHandlerFunc(GetResult))
	router.Handle("/results/user/{userKey}/challenge/{challengeKey}", ContextHandlerFunc(GetResultForUserChallenge))
}

// CreateResult saves a new result when a coder starts a challenge.
func CreateResult(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, nil
	}

	var body = struct {
		ChallengeKey string
	}{}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	var profiles model.Profiles
	keys, err := model.NewQueryForProfile().
		Ancestor(p.User).
		GetAll(ctx, &profiles)

	if len(keys) != 1 {
		return http.StatusInternalServerError, errors.New("Profile not found")
	}

	err = json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		return http.StatusBadRequest, err
	}

	key, err := datastore.DecodeKey(body.ChallengeKey)

	if err != nil {
		return http.StatusBadRequest, err
	}

	var results []model.Result

	resultKeys, err := model.NewQueryForResult().
		Ancestor(keys[0]).
		Filter("Challenge = ", key).
		Limit(1).
		GetAll(ctx, &results)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if len(resultKeys) == 1 {
		var u model.User
		if err = datastore.Get(ctx, p.User, &u); err != nil {
			return http.StatusInternalServerError, nil
		}

		if results[0].Finished.Equal(time.Time{}) || u.Company != nil {
			json.NewEncoder(w).Encode(results[0].Key(resultKeys[0]))
			return http.StatusOK, nil
		}

		return http.StatusForbidden, errors.New("you already finished this challenge")
	}

	var challenge model.Challenge
	if err = datastore.Get(ctx, key, &challenge); err != nil {
		return http.StatusInternalServerError, err
	}

	result := model.Result{
		Challenge:        key,
		StartTimes:       make([]time.Time, len(challenge.Tasks)),
		FinalSubmissions: make([]*datastore.Key, len(challenge.Tasks)),
		Started:          time.Now(),
	}
	key, err = result.PutWithParent(ctx, keys[0])
	if err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(result.Key(key))
	return http.StatusOK, nil
}

func GetResultForUserChallenge(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}
	userKey, err := datastore.DecodeKey(mux.Vars(r)["userKey"])
	if err != nil {
		return http.StatusBadRequest, err
	}
	challengeKey, err := datastore.DecodeKey(mux.Vars(r)["challengeKey"])
	if err != nil {
		return http.StatusBadRequest, err
	}

	keys, err := model.NewQueryForProfile().
		Ancestor(userKey).
		Limit(1).
		KeysOnly().
		GetAll(ctx, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if len(keys) != 1 {
		return http.StatusNotFound, nil
	}

	var results model.Results
	resultKeys, err := model.NewQueryForResult().
		Filter("Challenge =", challengeKey).
		Ancestor(keys[0]).
		Limit(1).
		GetAll(ctx, &results)

	if err != nil {
		return http.StatusInternalServerError, err
	}
	if len(resultKeys) != 1 {
		return http.StatusNotFound, nil
	}
	json.NewEncoder(w).Encode(results[0].Key(resultKeys[0]))
	return http.StatusOK, nil
}

func GetResult(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	resultKey, err := datastore.DecodeKey(mux.Vars(r)["resultKey"])
	if err != nil {
		return http.StatusBadRequest, err
	}

	var result model.Result
	if err := datastore.Get(ctx, resultKey, &result); err != nil {
		return http.StatusInternalServerError, nil
	}

	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	var u model.User
	if err = datastore.Get(ctx, p.User, &u); err != nil {
		return http.StatusInternalServerError, nil
	}

	if u.Company == nil && !util.HasParent(p.User, resultKey) {
		return http.StatusUnauthorized, nil
	}

	if result.Finished.Equal(time.Time{}) {
		if util.HasParent(p.User, resultKey) {
			return createFinalResult(ctx, w, *result.Key(resultKey))
		}
		var challenge model.Challenge
		if err := datastore.Get(ctx, result.Challenge, &challenge); err != nil {
			return http.StatusInternalServerError, err
		}
		if u.Company != nil && result.Started.Add(challenge.Duration).Before(time.Now()) {
			return createFinalResult(ctx, w, *result.Key(resultKey))
		}
	}

	json.NewEncoder(w).Encode(result.Key(resultKey))
	return http.StatusOK, nil
}

func createFinalResult(ctx context.Context, w http.ResponseWriter, result model.KeyedResult) (int, error) {
	var challenge model.Challenge
	if err := datastore.Get(ctx, result.Challenge, &challenge); err != nil {
		return http.StatusInternalServerError, nil
	}

	result.Finished = time.Now()

	for i, taskKey := range challenge.Tasks {
		key, err := getLatestSubmissionKey(ctx, result.Key, taskKey)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		result.FinalSubmissions[i] = key
	}
	_, err := result.Put(ctx, result.Key)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	json.NewEncoder(w).Encode(result)

	go computeFinalScore(ctx, result, challenge)

	return http.StatusOK, nil
}

func getLatestSubmissionKey(ctx context.Context, resultKey, taskKey *datastore.Key) (*datastore.Key, error) {
	keys, err := model.NewQueryForSubmission().
		Ancestor(resultKey).
		Filter("Task =", taskKey).
		Order("-Time").
		KeysOnly().
		Limit(1).
		GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}
	if len(keys) != 1 {
		return nil, nil
	}
	return keys[0], nil
}

func computeFinalScore(ctx context.Context, result model.KeyedResult, challenge model.Challenge) {
	if err := logic.Resulter(challenge.Resulter).Call(ctx, result, challenge); err != nil {
		log.Warningf(ctx, "resulter failed: %s", err.Error())
	}
}
