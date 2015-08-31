package model

//go:generate generator

type Test struct {
	Tester int
	Name   string
	Params map[string]string
}
