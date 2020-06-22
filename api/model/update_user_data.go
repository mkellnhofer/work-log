package model

// UpdateUser contains information about a user.
type UpdateUserData struct {
	Name     string              `json:"name"`
	Username string              `json:"username"`
	Contract *UpdateUserContract `json:"contract"`
}
