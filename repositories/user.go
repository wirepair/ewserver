package repositories

import (
	"github.com/wirepair/ewserver/errors"
	"github.com/wirepair/ewserver/store"
	"github.com/wirepair/ewserver/types"
)

const userTable = "users"

type UserRepository interface {
	AddUser(user *types.User) error
	DeleteUser(user *types.User) error
	UpdateUser(user *types.User) (bool, error)
	FindByID(ID int) (*types.User, error)
	FindByUserName(userName types.UserName) (*types.User, error)
	AccountDisabled(user *types.User) bool
}

type UserStore struct {
	store store.Storer
}

func NewUserStore(store store.Storer) *UserStore {
	return &UserStore{store: store}
}

func (userStore *UserStore) AddUser(user *types.User) error {
	exists, err := userStore.FindByID(user.ID)
	if err != nil {
		return err
	}

	if exists != nil {
		return errors.E(user.UserName, errors.Op("add"), errors.Exist)
	}

	return userStore.store.StoreUser(user)
}

func (userStore *UserStore) UpdateUser(user *types.User) error {
	exists, err := userStore.FindByID(user.ID)
	if err != nil {
		return err
	}

	if exists == nil {
		return errors.E(user.UserName, errors.Op("update"), errors.NotExist)
	}

	return userStore.store.StoreUser(user)
}

func (userStore *UserStore) FindByID(ID []byte) (*types.User, error) {
	return userStore.store.FindUserByID(ID)
}

func (userStore *UserStore) FindByUserName(userName types.UserName) (*types.User, error) {
	return userStore.store.FindUserByName(string(userName))
}
