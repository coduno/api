package models

// CodeData is the data to receive from the codeground
type CodeData struct {
	CodeBase string `json:"codeBase"`
	Token    string `json:"token"`
	Language string `json:"language"`
}
