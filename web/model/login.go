package model

// Login stores data for the login view.
type Login struct {
	ErrorMessage string
}

// NewLogin creates a new Login view model.
func NewLogin() *Login {
	return &Login{}
}
