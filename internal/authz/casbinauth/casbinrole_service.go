package casbinauth

import (
	"errors"
	"strings"

	"github.com/casbin/casbin"
)

var (
	// ErrAddUserToGroup when unable to add a user to a group
	ErrAddUserToGroup = errors.New("unable to add user to group")
	// ErrDeleteUserFromGroup when we can't remove a user from a group
	ErrDeleteUserFromGroup = errors.New("unable to remove user from group")
	// ErrAddPermission when we are unable to add a new permission
	ErrAddPermission = errors.New("unable to add permission")
	// ErrDeletePermission when we are unable to delete a permission
	ErrDeletePermission = errors.New("unable to delete permission")
	// ErrInvalidResource when the resource doesn't look like a URL
	ErrInvalidResource = errors.New("invalid resource specified, must start with /")
)

// CasbinRoleService implements role management via casbin
type CasbinRoleService struct {
	enforcer *casbin.SyncedEnforcer
}

// NewRoleService creates a new role service backed by casbin
func NewRoleService(enforcer *casbin.SyncedEnforcer) *CasbinRoleService {
	return &CasbinRoleService{enforcer: enforcer}
}

// Users returns all users
func (r *CasbinRoleService) Users() []string {
	return r.enforcer.GetAllNamedSubjects("p")
}

// Groups returns all groups
func (r *CasbinRoleService) Groups() []string {
	return r.enforcer.GetAllRoles()
}

// Permissions defined for this role service
func (r *CasbinRoleService) Permissions() [][]string {
	return r.enforcer.GetNamedPolicy("p")
}

// DeleteGroup from the service
func (r *CasbinRoleService) DeleteGroup(group string) error {
	r.enforcer.DeleteRole(group)
	return nil
}

// AddUserToGroup ...
func (r *CasbinRoleService) AddUserToGroup(user, group string) error {
	if ok := r.enforcer.AddRoleForUser(user, group); !ok {
		return ErrAddUserToGroup
	}
	return nil
}

// DeleteUserFromGroup ...
func (r *CasbinRoleService) DeleteUserFromGroup(user, group string) error {
	if ok := r.enforcer.DeleteRoleForUser(user, group); !ok {
		return ErrDeleteUserFromGroup
	}

	users := r.enforcer.GetUsersForRole(group)
	if users == nil || len(users) == 0 {
		return r.DeleteGroup(group)
	}
	return nil
}

// AddPermission adds an allow permission for either a user/group to access an object using the supplied method
func (r *CasbinRoleService) AddPermission(subject, object, method string) error {
	if !strings.HasPrefix(object, "/") {
		return ErrInvalidResource
	}

	if ok := r.enforcer.AddPolicy(subject, object, method); !ok {
		return ErrAddPermission
	}
	return nil
}

// DeletePermission for either a user/group to access an object using the supplied method
func (r *CasbinRoleService) DeletePermission(subject, object, method string) error {
	if ok := r.enforcer.DeletePermission(subject, object, method); !ok {
		return ErrDeletePermission
	}
	return nil
}
