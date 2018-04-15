package mock

import (
	"github.com/wirepair/ewserver/ewserver"
)

// UserService represents a mock implementation of myapp.UserService.
type UserService struct {
	InitFn      func() error
	InitInvoked bool

	UserFn      func(userName ewserver.UserName) (*ewserver.User, error)
	UserInvoked bool

	UsersFn      func() ([]*ewserver.User, error)
	UsersInvoked bool

	AuthenticateFn      func(userName ewserver.UserName, password string) (*ewserver.User, error)
	AuthenticateInvoked bool

	ChangePasswordFn      func(userName ewserver.UserName, current, new string) error
	ChangePasswordInvoked bool

	ResetPasswordFn      func(userName ewserver.UserName, new string) error
	ResetPasswordInvoked bool

	CreateFn      func(user *ewserver.User, password string) error
	CreateInvoked bool

	DeleteFn      func(userName ewserver.UserName) error
	DeleteInvoked bool
}

// User invokes the mock implementation and marks the function as invoked.

// Init the userBucket
func (u *UserService) Init() error {
	u.InitInvoked = true
	return u.InitFn()
}

// Authenticate a user to grant access, returns the User on success, error otherwise.
func (u *UserService) Authenticate(userName ewserver.UserName, password string) (*ewserver.User, error) {
	u.AuthenticateInvoked = true
	return u.AuthenticateFn(userName, password)
}

// ChangePassword of a user, provided they exist and the current password matches.
func (u *UserService) ChangePassword(userName ewserver.UserName, current, new string) error {
	u.ChangePasswordInvoked = true
	return u.ChangePasswordFn(userName, current, new)
}

// ResetPassword of a user, provided they exist
func (u *UserService) ResetPassword(userName ewserver.UserName, new string) error {
	u.ResetPasswordInvoked = true
	return u.ResetPasswordFn(userName, new)
}

// User finds the user by ID.
func (u *UserService) User(userName ewserver.UserName) (*ewserver.User, error) {
	u.UserInvoked = true
	return u.UserFn(userName)
}

// Users returns all users
func (u *UserService) Users() ([]*ewserver.User, error) {
	u.UsersInvoked = true
	return u.UsersFn()
}

// Create adds a new user if it does not already exist
func (u *UserService) Create(user *ewserver.User, password string) error {
	u.CreateInvoked = true
	return u.Create(user, password)
}

// Delete a User from the system. Does not return an error if user does not exist
func (u *UserService) Delete(userName ewserver.UserName) error {
	u.DeleteInvoked = true
	return u.Delete(userName)
}
