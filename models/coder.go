package models

// Coder contains the data related to a coder
type Coder struct {
	EntityID  int64  `datastore:"-" json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
