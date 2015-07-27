package models

import "time"

const AccessTokenKind = "accesstokens"

type AccessToken struct {
	Value        []byte    `json:"value"`
	Scopes       []string  `json:"scopes"`
	Description  string    `json:"description"`
	Creation     time.Time `json:"creation"`
	Modification time.Time `json:"modification"`
	Expiry       time.Time `json:"expiry"`
}
