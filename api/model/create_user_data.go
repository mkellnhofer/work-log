package model

// CreateUser
//
// Holds information about a new user.
//
// swagger:model CreateUser
type CreateUserData struct {
	// The name of the user.
	// min length: 1
	// max length: 100
	// example: John
	Name string `json:"name"`

	// The username of the user.
	// min length: 1
	// max length: 100
	// example: john
	Username string `json:"username"`

	// The password of the user.
	// min length: 1
	// max length: 100
	// example: secret
	Password string `json:"password"`

	// The work contract of the user.
	Contract *CreateContract `json:"contract"`
}
