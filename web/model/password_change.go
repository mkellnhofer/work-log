package model

// PasswordChange stores data for the password change view.
type PasswordChange struct {
	ErrorMessage string
}

// NewPasswordChange creates a new PasswordChange view model.
func NewPasswordChange() *PasswordChange {
	return &PasswordChange{}
}
