package model

type TestStats struct {
	Stdout,
	Stderr string
	Test   int64
	Failed bool
}
