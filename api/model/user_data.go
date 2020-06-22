package model

// UserData contains information about a user.
type UserData struct {
	Id       int           `json:"id"`
	Name     string        `json:"name"`
	Username string        `json:"username"`
	Contract *UserContract `json:"contract"`
}
