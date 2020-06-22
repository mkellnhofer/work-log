package model

// Error
//
// Supplies information about an error.
//
// swagger:model Error
type Error struct {
	// error code
	Code int `json:"code"`
	// error message
	Message string `json:"message"`
}

// NewError creates a new Error model.
func NewError(code int, message string) *Error {
	return &Error{code, message}
}
