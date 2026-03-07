package model

// TokenList
//
// A list of tokens.
//
// swagger:model TokenList
type TokenList struct {
	// The list of tokens.
	Tokens []*Token `json:"tokens"`
}

// NewTokenList creates a new TokenList model.
func NewTokenList(tokens []*Token) *TokenList {
	return &TokenList{tokens}
}
