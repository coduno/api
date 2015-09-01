package model

import (
	"bytes"
	"encoding/gob"

	"google.golang.org/appengine/datastore"
)

//go:generate generator

type Test struct {
	Tester int
	Name   string
	Params map[string]string
}

func (t *Test) Load(ps []datastore.Property) error {
	for _, p := range ps {
		switch p.Name {
		case "Tester":
			t.Tester = p.Value.(int)
		case "Name":
			t.Name = p.Value.(string)
		case "Params":
			if err := gob.NewDecoder(bytes.NewReader(p.Value.([]byte))).Decode(&t.Params); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Test) Save() ([]datastore.Property, error) {
	buf := new(bytes.Buffer)

	if err := gob.NewEncoder(buf).Encode(t.Params); err != nil {
		return nil, err
	}

	return []datastore.Property{
		{
			Name:  "Tester",
			Value: int64(t.Tester),
		},
		{
			Name:  "Name",
			Value: t.Name,
		},
		{
			Name:    "Params",
			Value:   buf.Bytes(),
			NoIndex: true,
		},
	}, nil
}
