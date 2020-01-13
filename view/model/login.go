package model

// Login stores login view data.
type Login struct {
	ErrorMessage string
}

// NewLogin creates a new Login view model.
func NewLogin() *Login {
	return &Login{}
}
