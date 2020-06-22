package model

// UpdateUserRoles contains information about the roles of a user.
type UpdateUserRoles struct {
	Roles []string `json:"roles"`
}
