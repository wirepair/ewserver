package ewserver

// RoleService for managing roles and permissions
type RoleService interface {
	RoleNames() []string                                   // lists role names
	RoleMap() [][]string                                   // lists subject to role mapping
	Permissions() [][]string                               // lists permissions for roles
	DeleteRole(roleName string) error                      // deletes all permissions related to this role
	AddSubjectToRole(subject, roleName string) error       // adds a subject to a role, creating the role if it does not exist
	DeleteSubjectFromRole(subject, roleName string) error  // deletes a subject from a role, if the only subject in the role, it deletes the role.
	AddPermission(subject, object, method string) error    // adds a new permission for a subject/role
	DeletePermission(subject, object, method string) error // deletes the permission
}
