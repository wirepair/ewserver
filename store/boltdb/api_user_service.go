package boltdb

import (
	"bytes"
	"encoding/base64"

	"github.com/boltdb/bolt"
	"github.com/wirepair/ewserver/ewserver"
)

const (
	apiKeyBucket = "api_keys"
	apiIDSize    = 16
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

// APIUserByID finds the user by their ID this is an O(N) operation, primarly used for admin management.
func (u *APIUserService) APIUserByID(ID []byte) (*ewserver.APIUser, error) {
	var foundUser *ewserver.APIUser

	id, err := base64.StdEncoding.DecodeString(string(ID))
	if err != nil {
		return nil, err
	}

	err = u.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(apiKeyBucket))
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			apiUser, err := ewserver.DecodeAPIUser(v)
			if err != nil {
				return err
			}

			if bytes.Equal(apiUser.ID, id) {
				foundUser = apiUser
				return nil
			}
		}
		return ewserver.ErrUserNotFound
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

// Create adds a new API key if it does not already exist, generates a random ID for API User management.
// We use the APIKey as the bucket key due to the majority of reads being done when requests contain
// the APIKey. If we used the ID, we would ineffiently be scanning the entire bucket for all API access.
func (u *APIUserService) Create(apiUser *ewserver.APIUser) error {
	var err error

	if exists, _ := u.APIUser(apiUser.Key); exists != nil {
		return ewserver.ErrUserAlreadyExists
	}

	apiUser.ID, err = ewserver.GenerateRandomBytes(apiIDSize)
	if err != nil {
		return err
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
func (u *APIUserService) Delete(apiKey ewserver.APIKey) error {
	return u.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(apiKeyBucket))
		return bucket.Delete(apiKey.Bytes())
	})
}
