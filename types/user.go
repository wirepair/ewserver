package types

import (
	"bytes"
	"encoding/gob"
)

// User represents a user with UI access
type User struct {
	ID        []byte
	FirstName string
	LastName  string
	UserName  UserName
	APIKey    APIKey
}

// NewUser from bytes
func NewUser() *User {
	return &User{}
}

// Encode the User into a gob of bytes
func (u *User) Encode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(u)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DecodeUser from bytes using gob and return a User.
func DecodeUser(userBytes []byte) (*User, error) {
	buf := bytes.NewBuffer(userBytes)
	enc := gob.NewDecoder(buf)
	u := NewUser()
	err := enc.Decode(u)
	return u, err
}
