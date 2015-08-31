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
		json.NewEncoder(w).Encode(results[0].Key(resultKeys[0]))
		return http.StatusOK, nil
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

// GetResultsByChallenge queries the results for a certain challenge to be reviewed by a company.
func GetResultsByChallenge(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, nil
	}

	key, err := datastore.DecodeKey(mux.Vars(r)["key"])
	if err != nil {
		return http.StatusBadRequest, err
	}

	var results model.Results
	keys, err := model.NewQueryForResult().
		Filter("Challenge =", key).
		GetAll(ctx, &results)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(results.Key(keys))
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

	if u.Company == nil && !util.HasParent(resultKey, p.User) {
		return http.StatusUnauthorized, nil
	}

	if result.Finished == (time.Time{}) {
		if util.HasParent(resultKey, p.User) {
			return createFinalResult(ctx, w, resultKey, result)
		}
		var challenge model.Challenge
		if err := datastore.Get(ctx, result.Challenge, &challenge); err != nil {
			return http.StatusInternalServerError, err
		}
		if u.Company != nil && result.Started.Add(challenge.Duration).Before(time.Now()) {
			return createFinalResult(ctx, w, resultKey, result)
		}
	}

	json.NewEncoder(w).Encode(result.Key(resultKey))
	return http.StatusOK, nil
}

func createFinalResult(ctx context.Context, w http.ResponseWriter, resultKey *datastore.Key, result model.Result) (int, error) {
	var challenge model.Challenge
	if err := datastore.Get(ctx, result.Challenge, &challenge); err != nil {
		return http.StatusInternalServerError, nil
	}

	result.Finished = time.Now()

	for i, taskKey := range challenge.Tasks {
		key, err := getLatestSubmissionKey(ctx, resultKey, taskKey)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		result.FinalSubmissions[i] = key
	}
	key, err := result.Put(ctx, resultKey)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	json.NewEncoder(w).Encode(result.Key(key))

	go computeFinalScore(ctx, challenge, result, resultKey)

	return http.StatusOK, nil
}

func getLatestSubmissionKey(ctx context.Context, resultKey, taskKey *datastore.Key) (*datastore.Key, error) {
	keys, err := datastore.NewQuery("").
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
		return nil, errors.New("found wrong number of final submissions")
	}
	return keys[0], nil
}

func computeFinalScore(ctx context.Context, challenge model.Challenge, result model.Result, resultKey *datastore.Key) {
	if err := logic.Resulter(challenge.Resulter).Call(ctx, resultKey); err != nil {
		log.Warningf(ctx, "resulter failed: %s", err.Error())
	}
}
