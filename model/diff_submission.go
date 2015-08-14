package model

//go:generate generator

// DiffSubmission holds the result of an outputtest.
type DiffSubmission struct {
	CodeSubmission

	DiffLines []int `datastore:",noindex"`
}
