package ewserver

import (
	"bytes"
	"encoding/gob"
)

// User represents a user with UI access
type User struct {
	UserName        UserName `json:"user_name"`
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	LastAddress     string   `json:"last_address"`     // Last IP Address that authenticated for this user
	RoleAssignments []string `json:"role_assignments"` // Roles this user is assigned to
	Password        []byte   `json:"-"`                // Becareful with this field.
}

// NewUser creates a new user
func NewUser() *User {
	u := &User{}
	return u
}

// Encode the User into a gob of bytes
func (u *User) Encode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(u); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DecodeUser from bytes using gob decoder and return a User.
func DecodeUser(userBytes []byte) (*User, error) {
	buf := bytes.NewBuffer(userBytes)
	enc := gob.NewDecoder(buf)
	u := NewUser()
	err := enc.Decode(u)
	return u, err
}

// UserService manages how users are accessed
type UserService interface {
	Init() error                                                        // Init the user service (prepare the tables/bucket whatever)
	Authenticate(userName UserName, password string) (*User, error)     // Authenticate the user with provided password
	Create(u *User, password string) error                              // Create the user with the supplied password
	ResetPassword(userName UserName, new string) error                  // Reset the user's password (should only be admin level)
	ChangePassword(userName UserName, current string, new string) error // ChangePassword for allowing users to change their password
	Delete(userName UserName) error                                     // Delete the user (admin only)
	User(userName UserName) (*User, error)                              // User returns the entire user
	Users() ([]*User, error)                                            // Users returns all users
}

// AuthnService allows a user to authenticate or change their password
type AuthnService interface {
	Authenticate(userName UserName, password string) (*User, error)     // Authenticate the user with provided password
	ChangePassword(userName UserName, current string, new string) error // ChangePassword for allowing users to change their password
	User(userName UserName) (*User, error)                              // User returns the entire user
}
