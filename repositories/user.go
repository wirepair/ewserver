package repositories

import (
	"github.com/wirepair/ewserver/errors"
	"github.com/wirepair/ewserver/store"
	"github.com/wirepair/ewserver/types"
)

// UserRepositorer interface for dealing only with user related
// queries
type UserRepositorer interface {
	AddUser(user *types.User) error
	DeleteUser(user *types.User) error
	UpdateUser(user *types.User) (bool, error)
	GetAllUsers() ([]*types.User, error)
	FindByID(ID int) (*types.User, error)
	FindByUserName(userName types.UserName) (*types.User, error)
}

// UserRepository implementation of the UserRepositorer separates concerns
// from accessing the data (store.Storer) from the logic necessary to
// query or process the results of a query. This way if you need to create
// a new service (like RPC) you don't have to add all the processing code
// to the RPC handlers, everything is contained here. Likewise, the database
// interface/implementation stays small because it is simply for getting
// or updating data
type UserRepository struct {
	store store.Storer
}

// NewUserRepository for processing user related requests to the data store
func NewUserRepository(store store.Storer) *UserRepository {
	return &UserRepository{store: store}
}

// AddUser to the user store
func (userRepo *UserRepository) AddUser(user *types.User) error {
	const op errors.Op = "UserRepository/AddUser"
	foundUser, err := userRepo.FindByUserName(user.UserName)
	if err != nil {
		return errors.E(user.UserName, op, err)
	}

	if foundUser != nil {
		return errors.E(user.UserName, op, errors.Exist)
	}

	return userRepo.store.StoreUser(user)
}

// UpdateUser details in the UserStore
func (userRepo *UserRepository) UpdateUser(user *types.User) error {
	const op errors.Op = "UserRepository/UpdateUser"
	foundUser, err := userRepo.FindByUserName(user.UserName)
	if err != nil {
		return errors.E(user.UserName, op, err)
	}

	if foundUser == nil {
		return errors.E(user.UserName, op, errors.NotExist)
	}

	return userRepo.store.StoreUser(user)
}

// DeleteUser from the UserStore
func (userRepo *UserRepository) DeleteUser(user *types.User) error {
	const op errors.Op = "UserRepository/DeleteUser"
	foundUser, err := userRepo.FindByUserName(user.UserName)
	if err != nil {
		return errors.E(user.UserName, op, err)
	}

	if foundUser == nil {
		return errors.E(user.UserName, op, errors.NotExist)
	}

	return userRepo.store.DeleteUserByName(user.UserName)
}

// FindByID returns the user given the provided ID.
func (userRepo *UserRepository) FindByID(ID []byte) (*types.User, error) {
	return userRepo.store.FindUserByID(ID)
}

// FindByUserName returns the user given the provided UserName
func (userRepo *UserRepository) FindByUserName(userName types.UserName) (*types.User, error) {
	return userRepo.store.FindUserByUserName(userName)
}

// FindByAPIKey returns the user given the provided APIKey
func (userRepo *UserRepository) FindByAPIKey(key types.APIKey) (*types.User, error) {
	return userRepo.store.FindUserByAPIKey(key.Bytes())
}

// FindAllUsers returns all users details.
func (userRepo *UserRepository) FindAllUsers() ([]*types.User, error) {
	return userRepo.store.FindAllUsers()
}
