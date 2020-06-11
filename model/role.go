package model

// Role defines the role of a user in the application.
type Role string

func (r Role) String() string {
	return string(r)
}

// Available roles.
const (
	RoleAdmin     Role = "admin"
	RoleEvaluator Role = "evaluator"
	RoleUser      Role = "user"
)

// Roles holds a list of roles.
var Roles = []Role{
	RoleAdmin,
	RoleEvaluator,
	RoleUser,
}

// Right defines the right to perform a action in the application.
type Right string

func (r Right) String() string {
	return string(r)
}

// Available rights.
const (
	RightGetUserData         Right = "get_user_data"
	RightChangeUserData      Right = "change_user_data"
	RightGetUserAccount      Right = "get_user_account"
	RightChangeUserAccount   Right = "change_user_account"
	RightGetEntryCharacts    Right = "get_entry_characteristics"
	RightChangeEntryCharacts Right = "change_entry_characteristics"
	RightGetAllEntries       Right = "get_all_entries"
	RightChangeAllEntries    Right = "change_all_entries"
	RightGetOwnEntries       Right = "get_own_entries"
	RightChangeOwnEntries    Right = "change_own_entries"
)

// RolesRights holds a mapping of roles and rights.
var RolesRights = map[Role][]Right{
	RoleAdmin:     roleAdminRights,
	RoleEvaluator: roleEvaluatorRights,
	RoleUser:      roleUserRights,
}

// Rights of the admin role.
var roleAdminRights = []Right{
	RightGetUserData,
	RightChangeUserData,
	RightGetUserAccount,
	RightChangeUserAccount,
	RightGetEntryCharacts,
	RightChangeEntryCharacts,
	RightGetAllEntries,
	RightChangeAllEntries,
}

// Rights of the evaluator role.
var roleEvaluatorRights = []Right{
	RightGetUserData,
	RightGetUserAccount,
	RightChangeUserAccount,
	RightGetEntryCharacts,
	RightGetAllEntries,
}

// Rights of the user role.
var roleUserRights = []Right{
	RightGetUserAccount,
	RightChangeUserAccount,
	RightGetEntryCharacts,
	RightGetOwnEntries,
	RightChangeOwnEntries,
}

// GetRoleRights returns the rights of a role.
func GetRoleRights(r Role) []Right {
	return RolesRights[r]
}
