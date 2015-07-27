package models

const UserKind = "users"

type User struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	HashedPassword []byte `json:"-"`
}
