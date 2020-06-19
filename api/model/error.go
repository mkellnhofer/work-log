package model

// Error stores error information.
type Error struct {
	code    int    `json:"code"`    // Error code
	message string `json:"message"` // Error message
}

// NewError creates a new Error model.
func NewError(code int, message string) *Error {
	return &Error{code, message}
}
