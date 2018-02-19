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
}

// NewUser from bytes
func NewUser() *User {
	return &User{}
}

func (u *User) Encode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(u)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecodedUser(userBytes []byte) (*User, error) {
	buf := bytes.NewBuffer(userBytes)
	enc := gob.NewDecoder(buf)
	u := NewUser()
	err := enc.Decode(u)
	return u, err
}
