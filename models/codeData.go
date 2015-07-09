package models

// CodeData is the data to receive from the codeground
type CodeData struct {
	EntityID string `json:"id"`
	CodeBase string `json:"codeBase"`
	Token    string `json:"token"`
}
