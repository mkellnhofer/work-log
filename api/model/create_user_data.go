package model

// CreateUserData contains information about a user.
type CreateUserData struct {
	Name     string              `json:"name"`
	Username string              `json:"username"`
	Password string              `json:"password"`
	Contract *CreateUserContract `json:"contract"`
}
