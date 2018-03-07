package ewserver

import (
	"bytes"
	"encoding/gob"
)

// User represents a user with UI access
type User struct {
	UserName    UserName
	FirstName   string
	LastName    string
	Password    []byte
	LastAddress string
	Roles       []*Role
}

// NewUser creates a new user
func NewUser() *User {
	u := &User{}
	return u
}

// Init the new user with a random ID and APIKey
func (u *User) Init() error {
	return nil
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
	Init() error
	Authenticate(userName UserName, password string) (*User, error)
	Create(u *User, password string) error
	ChangePassword(userName UserName, current string, new string) error
	Delete(userName UserName) error
	User(userName UserName) (*User, error)
	Users() ([]*User, error)
}
