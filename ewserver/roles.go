package ewserver

// RoleAction defines the types of actions a role can execute
type RoleAction uint8

// Create, Read, Update, Delete actions
const (
	Create RoleAction = iota
	Read
	Update
	Delete
)

// Role defines a simple Subject (user), Object (thing being accessed), and Action (CRUD)
type Role struct {
	Subject UserName
	Object  string
	Action  RoleAction
}
