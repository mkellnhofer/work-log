package model

// CreateToken
//
// Holds information for creating a new token.
//
// swagger:model CreateToken
type CreateToken struct {
	// The name of the token.
	// min length: 1
	// max length: 30
	// example: My API Token
	Name string `json:"name"`
}
