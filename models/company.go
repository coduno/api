package models

// Company contains the data related to a company
type Company struct {
	EntityID string `datastore:"-" json:"id"`
	Name     string `json:"name"`
}
