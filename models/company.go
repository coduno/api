package models

// Company -
type Company struct {
	EntityID int64  `datastore:"-" json:"id"`
	Name     string `json:"name"`
}
