package model

//go:generate generator

// CodeTask is any task where Users are asked to upload code.
type CodeTask struct {
	Task
	Flags     string   `datastore:",noindex"`
	Languages []string `datastore:",noindex"`
	Runner    string   `datastore:",noindex"`
}
