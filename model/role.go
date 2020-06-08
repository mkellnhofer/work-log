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
