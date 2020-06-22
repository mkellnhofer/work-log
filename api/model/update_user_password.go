package model

// UpdateUserPassword
//
// Holds the new password of a user.
//
// swagger:model UpdateUserPassword
type UpdateUserPassword struct {
	// The password.
	Password string `json:"password"`
}
