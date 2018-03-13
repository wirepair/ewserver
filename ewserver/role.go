package ewserver

import (
	"bytes"
	"encoding/gob"
)

// Role defines a simple Subject (user), Object (thing being accessed), and Action (CRUD)
type Role struct {
	Subject string `json:"subject"`
	Object  string `json:"object"`
	Action  string `json:"action"`
}

// NewRole creates a new role from the provided subject, object and action
func NewRole(subject, object, action string) *Role {
	return &Role{Subject: subject, Object: object, Action: action}
}

// Equals tests the equality of the individual properties of a role
func (r *Role) Equals(subject, object, action string) bool {
	if r.Subject == subject && r.Object == object && r.Action == action {
		return true
	}
	return false
}

// RolePermissions stores the object & actions of a Role
type RolePermissions struct {
	ObjectActions map[string][]string
}

// NewRolePermissions creates a new RolePermission to allow object & action
// assignment.
func NewRolePermissions() *RolePermissions {
	p := &RolePermissions{}
	p.ObjectActions = make(map[string][]string, 0)
	return p
}

// Add the allowed action to the object in our permissions, will create the object
// if it does not exist.
func (p *RolePermissions) Add(object, action string) {
	var actions []string
	var ok bool

	if actions, ok = p.ObjectActions[object]; !ok {
		actions = make([]string, 0)
	}

	actions = append(actions, action)
	p.ObjectActions[object] = actions
}

// DeleteObject removes it, and all allowed actions defined for it.
func (p *RolePermissions) DeleteObject(object string) {
	delete(p.ObjectActions, object)
}

// DeleteAction removes the action from the object.
func (p *RolePermissions) DeleteAction(object, action string) {
	var actions []string
	var ok bool

	if actions, ok = p.ObjectActions[object]; !ok {
		return
	}

	for i, value := range actions {
		if value == action {
			actions = append(actions[:i], actions[i+1:]...)
			break
		}
	}
}

// Encode the User into a gob of bytes
func (p *RolePermissions) Encode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(p); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DecodeRolePermissions from bytes using gob decoder and return RolePermissions.
func DecodeRolePermissions(permissionBytes []byte) (*RolePermissions, error) {
	buf := bytes.NewBuffer(permissionBytes)
	enc := gob.NewDecoder(buf)
	p := NewRolePermissions()
	err := enc.Decode(p)
	return p, err
}

// RoleService manages and enforces roles and authorization for resources
type RoleService interface {
	Create(role *Role) error
	Delete(role *Role) error
	Roles() ([]*Role, error)
	Authorize(subject string, object, action string) bool
}
