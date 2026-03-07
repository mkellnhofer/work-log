package model

// Token
//
// Contains information about a token.
//
// swagger:model Token
type Token struct {
	// The ID of the token.
	// example: 1
	Id int `json:"id"`

	// The name of the token.
	// min length: 1
	// max length: 30
	// example: My API Token
	Name string `json:"name"`

	// The token string.
	// example: a1b2c3d4...
	Token string `json:"token"`
}
