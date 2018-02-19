package types

import (
	"bytes"
	"encoding/gob"
)

// APIUser represents a user who only accesses via the API (UI).
type APIUser struct {
	ID  UserName
	Key string
}

// NewAPIUser from bytes
func NewAPIUser() *APIUser {
	return &APIUser{}
}

func (a *APIUser) Encode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(a)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecodedAPIUser(apiUserBytes []byte) (*APIUser, error) {
	buf := bytes.NewBuffer(apiUserBytes)
	enc := gob.NewDecoder(buf)
	a := NewAPIUser()
	err := enc.Decode(a)
	return a, err
}
