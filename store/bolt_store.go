package store

import (
	"bytes"
	"crypto/subtle"

	"github.com/boltdb/bolt"
	"github.com/wirepair/ewserver/types"
)

const (
	userBucket    = "users"    // bucket for storing UI users
	apiUserBucket = "apiUsers" // bucket for storing API users
)

// BoltStore for saving data to a Bolt DB file.
type BoltStore struct {
	config *Config
	db     *bolt.DB
}

// NewBoltStore for saving data to a bolt DB
func NewBoltStore() *BoltStore {
	b := &BoltStore{}
	return b
}

// Open the database file for writing
func (b *BoltStore) Open(config *Config) error {
	var err error

	b.db, err = bolt.Open(config.ConnectionString, 0600, nil)
	if err == nil {
		b.createBuckets()
	}
	return err
}

// createBuckets for ones that do not already exist.
func (b *BoltStore) createBuckets() error {
	if err := b.createBucket(userBucket); err != nil {
		return err
	}

	return b.createBucket(apiUserBucket)
}

// createBucket creates a new bucket with the provided name if it
// does not exist already.
func (b *BoltStore) createBucket(name string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		return err
	})
}

// StoreUser adds a new user or overrides the User.UserName key with
// the updated value.
func (b *BoltStore) StoreUser(user *types.User) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBucket))
		userBytes, err := user.Encode()
		if err != nil {
			return err
		}
		return bucket.Put(user.UserName.Bytes(), userBytes)
	})
}

// DeleteUserByName from the db file. Does not return an error if user does not exist
func (b *BoltStore) DeleteUserByName(userName types.UserName) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBucket))
		return bucket.Delete(userName.Bytes())
	})
}

// FindUserByUserName from the db file.
func (b *BoltStore) FindUserByUserName(userName types.UserName) (*types.User, error) {
	var foundUser *types.User

	err := b.db.View(func(tx *bolt.Tx) error {
		var decodeErr error
		bucket := tx.Bucket([]byte(userBucket))
		userBytes := bucket.Get(userName.Bytes())
		if userBytes == nil {
			return nil
		}

		foundUser, decodeErr = types.DecodeUser(userBytes)
		return decodeErr
	})
	return foundUser, err
}

// FindUserByID is an O(n) operation as it iterates over all users, decodes and
// compares the user.ID with the provided ID. FindByUserName is prefered since
// we store by UserName as the Key, not IDs.
func (b *BoltStore) FindUserByID(ID []byte) (*types.User, error) {
	var foundUser *types.User

	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBucket))
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			user, decodeErr := types.DecodeUser(v)
			if decodeErr != nil {
				continue
			}

			if bytes.Compare(user.ID, ID) == 0 {
				foundUser = user
				return nil
			}
		}
		return nil
	})
	return foundUser, err
}

// FindUserByAPIKey is an O(n) operation as it iterates over all users, decodes and
// constant time compares the user.Key against the provided Key.
func (b *BoltStore) FindUserByAPIKey(key []byte) (*types.User, error) {
	var foundUser *types.User

	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBucket))
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			user, decodeErr := types.DecodeUser(v)
			if decodeErr != nil {
				continue
			}

			// do a constant time to compare to prevent against brute force attacks
			if subtle.ConstantTimeCompare(key, user.APIKey.Bytes()) == 1 {
				// note we don't return immediately so this always takes exactly O(n)
				// in a real database we obviously wouldn't have this problem. But
				// it's not like you are using this in production on the internet right?
				foundUser = user
			}
		}
		return nil
	})
	return foundUser, err
}

// FindAllUsers returns a slice of all users
func (b *BoltStore) FindAllUsers() ([]*types.User, error) {
	users := make([]*types.User, 0)

	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBucket))
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			user, decodeErr := types.DecodeUser(v)
			if decodeErr != nil {
				continue
			}
			users = append(users, user)
		}
		return nil
	})
	return users, err
}

// Close the db file.
func (b *BoltStore) Close() error {
	if err := b.db.Sync(); err != nil {
		return err
	}
	return b.db.Close()
}
