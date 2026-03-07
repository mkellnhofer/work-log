package model

const (
	TokenLength = 32
)

// Token stores information about an API token.
type Token struct {
	Id        int
	UserId    int
	Name      string
	Token     string
	HashedToken string
	TruncatedToken string
}

// NewToken creates a new Token model with a generated token string.
func NewToken(userId int, name string) *Token {
	token := generateRandomString(TokenLength)
	hashedToken := createHashedString(token)
	truncatedToken := createTruncatedString(token, 4)
	return &Token{
		UserId:         userId,
		Name:           name,
		Token:          token,
		HashedToken:    hashedToken,
		TruncatedToken: truncatedToken,
	}
}

func IsValidTokenValue(tokenValue string) bool {
	return len(tokenValue) == TokenLength
}
