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

	if p.UserKey.Parent() != nil {
		json.NewEncoder(w).Encode(result.Key(resultKey))
		return http.StatusOK, nil
	}

	if !util.HasParent(resultKey, p.UserKey) {
		return http.StatusUnauthorized, nil
	}
	return createFinalResult(ctx, w, resultKey, result)
}

func createFinalResult(ctx context.Context, w http.ResponseWriter, resultKey *datastore.Key, result model.Result) (int, error) {
	go computeFinalScore(ctx, result)

	var challenge model.Challenge
	if err := datastore.Get(ctx, result.Challenge, &challenge); err != nil {
		return http.StatusInternalServerError, nil
	}

	var taskKey *datastore.Key
	for _, taskKey = range challenge.Tasks {
		switch taskKey.Kind() {
		case model.CodeTaskKind:
			var submissions model.CodeSubmissions
			keys, err := doQuery(ctx, model.NewQueryForCodeSubmission(), resultKey, taskKey, submissions)
			if err != nil {
				return http.StatusInternalServerError, nil
			}
			if len(keys) == 0 {
				// Most likely the authenticated user called this endpoint
				// before finishing the challenge
				return http.StatusUnauthorized, nil
			}
			result.FinalSubmissions = append(result.FinalSubmissions, keys[0])
		//TODO(pbochis, vbalan, flowlo): Add more cases when more task kinds are added.
		default:
			return http.StatusBadRequest, errors.New("Unknown submission kind.")
		}
	}
	key, err := result.Save(ctx, resultKey)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	json.NewEncoder(w).Encode(result.Key(key))
	return http.StatusOK, nil
}

func doQuery(ctx context.Context, query *datastore.Query, resultKey, taskKey *datastore.Key, dst interface{}) ([]*datastore.Key, error) {
	keys, err := query.
		Ancestor(resultKey).
		Filter("Task=", taskKey).
		Order("-Time").
		Limit(1).
		GetAll(ctx, &dst)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func computeFinalScore(ctx context.Context, result model.Result) {
	// Note: See comment above Logic in model/challenge.go. This can only be
	// calculated after challenge.Logic is clearly defined
}
