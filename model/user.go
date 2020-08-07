package model

// Standard user IDs.
const (
	SystemUserId    int = -1
	AnonymousUserId int = 0
)

// User stores information about a user.
type User struct {
	Id                 int    // ID of the user
	Name               string // Name of the user
	Username           string // Username of the user
	Password           string // Password of the user
	MustChangePassword bool   // Determines if user must change password
}

// NewUser creates a new User model.
func NewUser() *User {
	return &User{}
}
