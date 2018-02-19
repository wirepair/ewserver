package store

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/wirepair/ewserver/errors"
	"github.com/wirepair/ewserver/types"
)

const (
	userBucket    = "users"
	apiUserBucket = "apiUsers"
)

type KeyNotExistErr struct {
	Key string
}

func (k KeyNotExistErr) Error() string {
	return k.Key + "does not exist"
}

type BoltStore struct {
	config *Config
	db     *bolt.DB
}

func NewBoltStore() *BoltStore {
	b := &BoltStore{}
	b.requestMap = make(map[string]string)
	return b
}

func (b *BoltStore) Open(config *Config) error {
	db, err := bolt.Open(config.ConnectionString, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (b *BoltStore) createBuckets() error {
	if err := b.createBucket(userBucket); err != nil {
		return err
	}

	return b.createBucket(apiUserBucket)
}

func (b *BoltStore) createBucket(name string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		return err
	})
}

func (b *BoltStore) StoreUser(user *types.User) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))
		userBytes, err := user.Encode()
		if err != nil {
			return errors.E(user.UserName, errors.Op("store"), err)
		}
		return b.Put(user.ID, userBytes)
	})
}

func (b *BoltStore) DeleteUserByID(ID []byte) error {
	return nil
}

func (b *BoltStore) FindUserByID(ID []byte) (*types.User, error) {
	return nil, nil
}

func (b *BoltStore) StoreAPIUser(apiUser *types.APIUser) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(apiUserBucket))
		apiUserBytes, err := apiUser.Encode()
		if err != nil {
			return errors.E(apiUser.ID, errors.Op("store"), err)
		}
		return b.Put([]byte(apiUser.ID), apiUserBytes)
	})
}

func (b *BoltStore) DeleteAPIUserByID(ID []byte) error {

}

func (b *BoltStore) Close() error {
	return b.db.Close()
}
