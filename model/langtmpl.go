package model

import (
	"reflect"

	"google.golang.org/appengine/datastore"
)

type LanguageTemplates map[string][]StoredObject

func (m LanguageTemplates) Load(ps []datastore.Property) error {
	for _, p := range ps {
		os, ok := p.Value.([]StoredObject)
		if !ok {
			return &datastore.ErrFieldMismatch{
				FieldName:  p.Name,
				StructType: reflect.TypeOf(m),
				Reason:     "cannot assert type",
			}
		}
		m[p.Name] = os
	}
	return nil
}

func (m LanguageTemplates) Save() ([]datastore.Property, error) {
	ps := make([]datastore.Property, 0, len(m))
	for k, v := range m {
		ps = append(ps, datastore.Property{
			Name:  k,
			Value: v,
		})
	}
	return ps, nil
}
