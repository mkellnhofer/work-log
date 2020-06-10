package model

// UserData holds all information about a user.
// (This model is needed to get/create/change all user related data at once.)
type UserData struct {
	Id           int           // ID of the user
	User         *User         // User
	UserContract *UserContract // User contract
}

// NewUserData creates a new UserData model.
func NewUserData(id int, u *User, uc *UserContract) *UserData {
	return &UserData{id, u, uc}
}
