package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/wirepair/ewserver/ewserver"
)

const (
	apiKeyBucket = "api_keys"
)

// APIUserService implementation that manages access to Users
type APIUserService struct {
	DB *bolt.DB
}

// NewAPIUserService creates a new API user service backed by an already open boltdb
func NewAPIUserService(db *bolt.DB) *APIUserService {
	u := &APIUserService{DB: db}
	return u
}

// Init the API key bucket
func (u *APIUserService) Init() error {
	return u.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(apiKeyBucket))
		return err
	})
}

// APIUser finds the user by their APIKey.
func (u *APIUserService) APIUser(apiKey ewserver.APIKey) (*ewserver.APIUser, error) {
	var foundUser *ewserver.APIUser

	err := u.DB.View(func(tx *bolt.Tx) error {
		var decodeErr error

		bucket := tx.Bucket([]byte(apiKeyBucket))
		apiUserBytes := bucket.Get(apiKey.Bytes())
		if apiUserBytes == nil {
			return ewserver.ErrUserNotFound
		}
		foundUser, decodeErr = ewserver.DecodeAPIUser(apiUserBytes)
		return decodeErr
	})
	return foundUser, err
}

// APIUsers returns all API Users
func (u *APIUserService) APIUsers() ([]*ewserver.APIUser, error) {
	foundAPIUsers := make([]*ewserver.APIUser, 0)

	err := u.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(apiKeyBucket))
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			apiUser, err := ewserver.DecodeAPIUser(v)
			if err != nil {
				return err
			}

			foundAPIUsers = append(foundAPIUsers, apiUser)
		}
		return nil
	})

	return foundAPIUsers, err
}

// Create adds a new user if it does not already exist
func (u *APIUserService) Create(apiUser *ewserver.APIUser) error {
	if exists, _ := u.APIUser(apiUser.Key); exists != nil {
		return ewserver.ErrUserAlreadyExists
	}

	return u.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(apiKeyBucket))
		userBytes, err := apiUser.Encode()
		if err != nil {
			return err
		}
		return bucket.Put(apiUser.Key.Bytes(), userBytes)
	})
}

// Delete a User from the system. Does not return an error if user does not exist
func (u *APIUserService) Delete(apiKey ewserver.UserName) error {
	return u.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(apiKeyBucket))
		return bucket.Delete(apiKey.Bytes())
	})
}
