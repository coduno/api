package model

//go:generate generator

type CodeTask struct {
	Task
	Flags     string
	Languages []string
	Runner    string
}
