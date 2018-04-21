package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/wirepair/ewserver/ewserver"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 10      // bcrypt cost to use (change if necessary, default from bcrypt is 10)
	userBucket = "users" // bucket for storing UI users
)

// UserService implementation that manages access to Users
type UserService struct {
	DB *bolt.DB
}

// NewUserService creates a new user service backed by an already open boltdb
func NewUserService(db *bolt.DB) *UserService {
	u := &UserService{DB: db}
	return u
}

// Init the userBucket
func (u *UserService) Init() error {
	return u.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(userBucket))
		return err
	})
}

// Authenticate a user to grant access, returns the User on success, error otherwise.
func (u *UserService) Authenticate(userName ewserver.UserName, password string) (*ewserver.User, error) {
	validUser, err := u.User(userName)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(validUser.Password, []byte(password)); err != nil {
		return nil, ewserver.ErrInvalidPassword
	}

	return validUser, nil
}

// ChangePassword of a user, provided they exist and the current password matches.
func (u *UserService) ChangePassword(userName ewserver.UserName, current, new string) error {
	var validUser *ewserver.User

	// Do everything in a writable transaction so there's no oddness with getting users while updating.
	tx, err := u.DB.Begin(true)
	if err != nil {
		return err
	}

	bucket := tx.Bucket([]byte(userBucket))

	userBytes := bucket.Get(userName.Bytes())
	if userBytes == nil {
		tx.Rollback()
		return ewserver.ErrUserNotFound
	}

	validUser, err = ewserver.DecodeUser(userBytes)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := bcrypt.CompareHashAndPassword(validUser.Password, []byte(current)); err != nil {
		tx.Rollback()
		return ewserver.ErrInvalidPassword
	}

	if validUser.Password, err = u.hashPassword([]byte(new)); err != nil {
		tx.Rollback()
		return err
	}

	encodedUser, err := validUser.Encode()
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := bucket.Put(validUser.UserName.Bytes(), encodedUser); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// ResetPassword of a user, provided they exist
func (u *UserService) ResetPassword(userName ewserver.UserName, new string) error {
	var validUser *ewserver.User

	// Do everything in a writable transaction so there's no oddness with getting users while updating.
	tx, err := u.DB.Begin(true)
	if err != nil {
		return err
	}

	bucket := tx.Bucket([]byte(userBucket))

	userBytes := bucket.Get(userName.Bytes())
	if userBytes == nil {
		tx.Rollback()
		return ewserver.ErrUserNotFound
	}

	validUser, err = ewserver.DecodeUser(userBytes)
	if err != nil {
		tx.Rollback()
		return err
	}

	if validUser.Password, err = u.hashPassword([]byte(new)); err != nil {
		tx.Rollback()
		return err
	}

	encodedUser, err := validUser.Encode()
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := bucket.Put(validUser.UserName.Bytes(), encodedUser); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// User finds the user by ID.
func (u *UserService) User(userName ewserver.UserName) (*ewserver.User, error) {
	var foundUser *ewserver.User

	err := u.DB.View(func(tx *bolt.Tx) error {
		var decodeErr error
		bucket := tx.Bucket([]byte(userBucket))
		userBytes := bucket.Get(userName.Bytes())
		if userBytes == nil {
			return ewserver.ErrUserNotFound
		}
		foundUser, decodeErr = ewserver.DecodeUser(userBytes)
		return decodeErr
	})
	return foundUser, err
}

// Users returns all users
func (u *UserService) Users() ([]*ewserver.User, error) {
	foundUsers := make([]*ewserver.User, 0)

	err := u.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBucket))
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			user, err := ewserver.DecodeUser(v)
			if err != nil {
				return err
			}

			foundUsers = append(foundUsers, user)
		}
		return nil
	})
	return foundUsers, err
}

// Create adds a new user if it does not already exist
func (u *UserService) Create(user *ewserver.User, password string) error {
	var err error

	if exists, _ := u.User(user.UserName); exists != nil {
		return ewserver.ErrUserAlreadyExists
	}

	if u.invalid(user) {
		return ewserver.ErrInvalidUser
	}
	user.Password, err = u.hashPassword([]byte(password))
	if err != nil {
		return err
	}

	return u.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBucket))
		userBytes, err := user.Encode()
		if err != nil {
			return err
		}
		return bucket.Put(user.UserName.Bytes(), userBytes)
	})
}

// Update the user details
func (u *UserService) Update(user *ewserver.User) error {
	if exists, _ := u.User(user.UserName); exists == nil {
		return ewserver.ErrUserAlreadyExists
	}

	if u.invalid(user) {
		return ewserver.ErrInvalidUser
	}

	return u.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBucket))
		userBytes, err := user.Encode()
		if err != nil {
			return err
		}
		return bucket.Put(user.UserName.Bytes(), userBytes)
	})
}

// Delete a User from the system. Does not return an error if user does not exist
func (u *UserService) Delete(userName ewserver.UserName) error {
	return u.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBucket))
		return bucket.Delete(userName.Bytes())
	})
}

func (u *UserService) hashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcryptCost)
}

// invalid checks if the user or user fields are invalid. Implement any
// additional validation you'd like here.
func (u *UserService) invalid(user *ewserver.User) bool {
	if user.UserName == "" {
		return true
	}

	return false
}
