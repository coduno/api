package gitlab

import (
	"encoding/json"
)

type Author struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Commit struct {
	Author    Author `json:"author"`
	ID        string `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	URL       string `json:"url"`
}

type Repository struct {
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
	Name        string `json:"name"`
	URL         string `json:"url"`
}

type Push struct {
	After             string     `json:"after"`
	Before            string     `json:"before"`
	Commits           []Commit   `json:"commits"`
	ProjectID         float64    `json:"project_id"`
	Ref               string     `json:"ref"`
	Repository        Repository `json:"repository"`
	TotalCommitsCount float64    `json:"total_commits_count"`
	UserID            float64    `json:"user_id"`
	UserName          string     `json:"user_name"`
}

func NewPush(data []byte) (Push, error) {
	var push Push
	if err := json.Unmarshal(data, &push); err != nil {
		return push, err
	}

	return push, nil
}
