package mapper

import (
	am "kellnhofer.com/work-log/api/model"
	m "kellnhofer.com/work-log/pkg/model"
)

// ToTokens converts a list of logic token models to an API token list (truncated).
func ToTokens(tokens []*m.Token) *am.TokenList {
	if tokens == nil {
		return nil
	}

	items := make([]*am.Token, len(tokens))
	for i, t := range tokens {
		items[i] = ToToken(t)
	}

	return am.NewTokenList(items)
}

// ToToken converts a logic token model to an API token model (truncated token).
func ToToken(t *m.Token) *am.Token {
	if t == nil {
		return nil
	}

	var out am.Token
	out.Id = t.Id
	out.Name = t.Name
	out.Token = t.TruncatedToken
	return &out
}

// ToTokenFull converts a logic token model to an API token model (full token).
func ToTokenFull(t *m.Token) *am.Token {
	if t == nil {
		return nil
	}

	var out am.Token
	out.Id = t.Id
	out.Name = t.Name
	out.Token = t.Token
	return &out
}
