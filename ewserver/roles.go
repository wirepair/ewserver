package ewserver

// RoleService for managing roles and permissions
type RoleService interface {
	Users() []string                                       // list users
	Groups() []string                                      // list groups
	Permissions() [][]string                               // list permissions
	DeleteGroup(group string) error                        // deletes all permissions related to this group
	AddUserToGroup(user, group string) error               // adds a user to a group, creating the group if it does not exist
	DeleteUserFromGroup(user, group string) error          // deletes a user from a group, if the only user in the group, it deletes the group.
	AddPermission(subject, object, method string) error    // adds a new permission for a user or group
	DeletePermission(subject, object, method string) error // deletes the permission
}
