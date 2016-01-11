package model

import "time"

//go:generate generator

// Submission is a form of result for some Task.
//
// TODO(flowlo): As soon as we also store other submissions, implement a
// PropertyLoadSaver similar to this:
//
//	func (s *Submission) Load(ps []datastore.Property) error {
//		return datastore.LoadStruct(s, ps)
//	}
//
//	func (s *Submission) Save() ([]datastore.Property, error) {
//		if s.Code.Name != "" && s.Answers != nil {
//			return nil, errors.New("cannot save Code and Answers in one Submission")
//		}
//		return ...
//	}
//
type Submission struct {
	ID       int64
	Time     time.Time
	Task     int64
	Code     StoredObject
	Language string
}
