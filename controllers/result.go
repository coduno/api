package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"github.com/coduno/api/util/passenger"
	"github.com/gorilla/mux"
	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
)

// CreateResult saves a new result when a coder starts a challenge.
func CreateResult(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if !util.CheckMethod(r, "POST") {
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
		Ancestor(p.UserKey).
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
	key, err = result.SaveWithParent(ctx, keys[0])
	if err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(result.Key(key))
	return http.StatusOK, nil
}

// GetResultsByChallenge queries the results for a certain challenge to be reviewed by a company.
func GetResultsByChallenge(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if !util.CheckMethod(r, "GET") {
		return http.StatusMethodNotAllowed, nil
	}

	key, err := datastore.DecodeKey(mux.Vars(r)["key"])

	if err != nil {
		return http.StatusBadRequest, err
	}

	var results model.Results
	keys, err := model.NewQueryForResult().
		Filter("Challenge=", key).
		GetAll(ctx, &results)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(results.Key(keys))
	return http.StatusOK, nil
}

func GetResult(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if !util.CheckMethod(r, "GET") {
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

	if p.UserKey.Parent() != nil && !util.HasParent(resultKey, p.UserKey) {
		return http.StatusUnauthorized, nil
	}

	if result.Finished == (time.Time{}) {
		var challenge model.Challenge
		if err := datastore.Get(ctx, result.Challenge, &challenge); err != nil {
			return http.StatusInternalServerError, err
		}
		if p.UserKey.Parent() != nil && result.Started.Add(challenge.Duration).After(time.Now()) {
			json.NewEncoder(w).Encode(result.Key(resultKey))
			return http.StatusOK, nil
		}
		return createFinalResult(ctx, w, resultKey, result, challenge)
	}

	json.NewEncoder(w).Encode(result.Key(resultKey))
	return http.StatusOK, nil
}

func createFinalResult(ctx context.Context, w http.ResponseWriter, resultKey *datastore.Key, result model.Result, challenge model.Challenge) (status int, err error) {
	go computeFinalScore(ctx, result)

	result.Finished = time.Now()

	for i, taskKey := range challenge.Tasks {
		var key *datastore.Key
		if key, err = getLatestSubmissionKey(ctx, resultKey, taskKey); err != nil {
			return http.StatusInternalServerError, err
		}
		result.FinalSubmissions[i] = key
	}

	if _, err = result.Save(ctx, resultKey); err != nil {
		return http.StatusInternalServerError, err
	}
	json.NewEncoder(w).Encode(result.Key(resultKey))
	return
}

func getLatestSubmissionKey(ctx context.Context, resultKey, taskKey *datastore.Key) (*datastore.Key, error) {
	keys, err := datastore.NewQuery("").
		Ancestor(resultKey).
		Filter("Task=", taskKey).
		Order("-Time").
		KeysOnly().
		Limit(1).
		GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, errors.New("no submission found")
	}
	return keys[0], nil
}

func computeFinalScore(ctx context.Context, result model.Result) {
	// Note: See comment above Logic in model/challenge.go. This can only be
	// calculated after challenge.Logic is clearly defined
}
