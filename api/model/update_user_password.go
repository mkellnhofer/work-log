package model

// UpdateUserPassword contains the new password of a user.
type UpdateUserPassword struct {
	Password string `json:"password"`
}
