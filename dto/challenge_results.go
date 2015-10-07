package dto

import (
	"time"

	"github.com/coduno/api/model"

	"google.golang.org/appengine/datastore"
)

type ChallengeResults []ChallengeResult

func (r ChallengeResults) Len() int {
	return len(r)
}
func (r ChallengeResults) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
func (r ChallengeResults) Less(i, j int) bool {
	if r[i].TotalTime == 0 {
		return false
	}
	if r[j].TotalTime == 0 {
		return true
	}
	return r[i].TotalTime < r[j].TotalTime
}

type ChallengeResult struct {
	User model.KeyedUser
	// Index is important
	TaskResults []TaskResult
	TotalTime   time.Duration
}

type TaskResult struct {
	Task            *datastore.Key
	NrOfSubmissions int
	CodingTime      time.Duration
}
