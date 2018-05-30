package casbinauth

import (
	"errors"
	"strings"

	"github.com/casbin/casbin"
)

var (
	// ErrAddSubjectToRole when unable to add a user to a role
	ErrAddSubjectToRole = errors.New("unable to add subject to role")
	// ErrDeleteSubjectFromRole when we can't remove a subject from a role
	ErrDeleteSubjectFromRole = errors.New("unable to remove subject from role")
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

// RoleNames returns all named subjects
func (r *CasbinRoleService) RoleNames() []string {
	return r.enforcer.GetAllNamedSubjects("p")
}

// RoleMap returns all role -> subject mappings
func (r *CasbinRoleService) RoleMap() [][]string {
	return r.enforcer.GetNamedGroupingPolicy("g")
}

// Permissions returns all defined for this role service
func (r *CasbinRoleService) Permissions() [][]string {
	return r.enforcer.GetNamedPolicy("p")
}

// DeleteRole from the service
func (r *CasbinRoleService) DeleteRole(roleName string) error {
	r.enforcer.DeleteRole(roleName)
	return nil
}

// AddSubjectToRole adds a subject to a role
func (r *CasbinRoleService) AddSubjectToRole(subject, roleName string) error {
	if ok := r.enforcer.AddRoleForUser(subject, roleName); !ok {
		return ErrAddSubjectToRole
	}
	return nil
}

// DeleteSubjectFromRole removes the subject from the role
func (r *CasbinRoleService) DeleteSubjectFromRole(subject, roleName string) error {
	if ok := r.enforcer.DeleteRoleForUser(subject, roleName); !ok {
		return ErrDeleteSubjectFromRole
	}

	subjects := r.enforcer.GetUsersForRole(roleName)
	if subjects == nil || len(subjects) == 0 {
		return r.DeleteRole(roleName)
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
