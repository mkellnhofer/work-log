package model

// UserData holds all information about a user.
// (This model is needed to get/create/change all user related data at once.)
type UserData struct {
	Id       int       // ID of the user
	User     *User     // User
	Contract *Contract // Contract
}

// NewUserData creates a new UserData model.
func NewUserData(id int, u *User, c *Contract) *UserData {
	return &UserData{id, u, c}
}
