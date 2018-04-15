package mock

import "github.com/wirepair/ewserver/ewserver"

// APIUserService represents a mock implementation of ewserver.APIUserService.
type APIUserService struct {
	CreateFn      func(u *ewserver.APIUser) error
	CreateInvoked bool

	APIUserFn      func(Key ewserver.APIKey) (*ewserver.APIUser, error)
	APIUserInvoked bool

	APIUserByIDFn      func(ID []byte) (*ewserver.APIUser, error)
	APIUserByIDInvoked bool

	APIUsersFn      func() ([]*ewserver.APIUser, error)
	APIUsersInvoked bool

	DeleteFn      func(Key ewserver.APIKey) error
	DeleteInvoked bool
}

// APIUser finds the user by their APIKey.
func (u *APIUserService) APIUser(apiKey ewserver.APIKey) (*ewserver.APIUser, error) {
	u.APIUserInvoked = true
	return u.APIUserFn(apiKey)
}

// APIUserByID finds the user by their ID this is an O(N) operation, primarly used for admin management.
func (u *APIUserService) APIUserByID(ID []byte) (*ewserver.APIUser, error) {
	u.APIUserByIDInvoked = true
	return u.APIUserByIDFn(ID)
}

// APIUsers returns all API Users
func (u *APIUserService) APIUsers() ([]*ewserver.APIUser, error) {
	u.APIUsersInvoked = true
	return u.APIUsersFn()
}

// Create adds a new API key if it does not already exist
func (u *APIUserService) Create(apiUser *ewserver.APIUser) error {
	u.CreateInvoked = true
	return u.CreateFn(apiUser)
}

// Delete a User from the system. Does not return an error if user does not exist
func (u *APIUserService) Delete(apiKey ewserver.APIKey) error {
	u.DeleteInvoked = true
	return u.DeleteFn(apiKey)
}
