package model

import "time"

type SimpleTestResult struct {
	Stdout  string `datastore:",noindex"`
	Stderr  string `datastore:",noindex"`
	Exit    string `datastore:",noindex"`
	Prepare string `datastore:",noindex"`

	Rusage Rusage    `datastore:",noindex"`
	Start  time.Time `datastore:",index"`
	End    time.Time `datastore:",index"`
}
