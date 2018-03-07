package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/wirepair/ewserver/store"
)

// BoltStore for saving data to a Bolt DB file.
type BoltStore struct {
	db *bolt.DB
}

// NewBoltStore for saving data to a bolt DB
func NewBoltStore() *BoltStore {
	b := &BoltStore{}
	return b
}

// Open the database file for writing
func (b *BoltStore) Open(config *store.Config) error {
	var err error

	b.db, err = bolt.Open(config.Options["database"], 0600, nil)
	return err
}

// DB exposes the database to allow assignment in services.
func (b *BoltStore) DB() *bolt.DB {
	return b.db
}

// Close the DB file.
func (b *BoltStore) Close() error {
	if err := b.db.Sync(); err != nil {
		return err
	}
	return b.db.Close()
}
