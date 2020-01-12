package model

// Error stores error view data.
type Error struct {
	ErrorMessage string
}

// NewError creates a new Error view model.
func NewError(errorMessage string) *Error {
	return &Error{errorMessage}
}
