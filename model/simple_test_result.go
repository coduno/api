package model

import "time"

type SimpleTestResult struct {
	Stdout  string `datastore:",noindex",json:",omitempty"`
	Stderr  string `datastore:",noindex",json:",omitempty"`
	Exit    string `datastore:",noindex",json:",omitempty"`
	Prepare string `datastore:",noindex",json:",omitempty"`

	Rusage Rusage    `datastore:",noindex",json:",omitempty"`
	Start  time.Time `datastore:",index",json:",omitempty"`
	End    time.Time `datastore:",index",json:",omitempty"`
}
