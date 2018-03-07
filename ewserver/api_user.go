package ewserver

import (
	"bytes"
	"encoding/gob"
)

// APIUser represents an api user
type APIUser struct {
	Key         APIKey
	LastAddress string
	Roles       []*Role
}

// NewAPIUser from bytes
func NewAPIUser() *APIUser {
	return &APIUser{}
}

// Encode encodes the APIUser to a slice of bytes
func (a *APIUser) Encode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(a); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DecodeAPIUser from a slice of bytes to an APIUser
func DecodeAPIUser(apiUserBytes []byte) (*APIUser, error) {
	buf := bytes.NewBuffer(apiUserBytes)
	enc := gob.NewDecoder(buf)
	a := NewAPIUser()
	err := enc.Decode(a)
	return a, err
}

// APIUserService manages how API users are managed
type APIUserService interface {
	Create(u *APIUser) error
	APIUser(Key APIKey) (*APIUser, error)
	APIUsers() ([]*APIUser, error)
	Delete(Key APIKey) error
}