package model

import (
	"bytes"
	"encoding/gob"
	"strings"

	"google.golang.org/appengine/datastore"
)

//go:generate generator

// Task is a concrete piece of work that cannot
// be split any further.
//
// This type is very general and can be implemented in vrious
// ways, accordingly implementing logic to make this Task comparable
// to others with respect to it's SkillWeights.
type Task struct {
	// Returns details on the assignment that is covered by this task.
	Assignment Assignment

	// Says what skills are needed/exercised to complete
	// the Task.
	SkillWeights SkillWeights

	// Refers to some logic that looks at the Submissions
	// of this task and produces a set of skills that
	// represent how well the user did in doing this Task.
	// It is to be weighted by SkillWeights.
	Tasker    int64                     `json:"-"`
	Templates map[string][]StoredObject `json:"-"`
	Languages []string                  `json:",omitempty"`
}

func (t *Task) Load(ps []datastore.Property) error {
	err := datastore.LoadStruct(&t.Assignment, ps)
	if err != nil {
		if _, ok := err.(*datastore.ErrFieldMismatch); !ok {
			return err
		}
	}
	err = datastore.LoadStruct(&t.SkillWeights, ps)
	if err != nil {
		if _, ok := err.(*datastore.ErrFieldMismatch); !ok {
			return err
		}
	}

	for _, p := range ps {
		switch p.Name {
		case "Templates":
			if err := gob.NewDecoder(bytes.NewReader(p.Value.([]byte))).Decode(&t.Templates); err != nil {
				return err
			}
		case "Tasker":
			t.Tasker = p.Value.(int64)
		case "Languages":
			t.Languages = strings.Split(p.Value.(string), " ")
		}
	}
	return nil
}

func (t *Task) Save() ([]datastore.Property, error) {
	var ps []datastore.Property

	tps, err := datastore.SaveStruct(&t.Assignment)
	if err != nil {
		return nil, err
	}
	ps = append(ps, tps...)

	tps, err = datastore.SaveStruct(&t.SkillWeights)
	if err != nil {
		return nil, err
	}
	ps = append(ps, tps...)

	buf := new(bytes.Buffer)

	if err := gob.NewEncoder(buf).Encode(t.Templates); err != nil {
		return nil, err
	}

	return append(ps, []datastore.Property{
		{
			Name:    "Tasker",
			Value:   int64(t.Tasker),
			NoIndex: true,
		},
		{
			Name:  "Languages",
			Value: strings.Join(t.Languages, " "),
		},
		{
			Name:    "Templates",
			Value:   buf.Bytes(),
			NoIndex: true,
		},
	}...), nil
}
