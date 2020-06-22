package model

// Error stores error information.
type Error struct {
	Code    int    `json:"code"`    // Error code
	Message string `json:"message"` // Error message
}

// NewError creates a new Error model.
func NewError(code int, message string) *Error {
	return &Error{code, message}
}
