package store

import "github.com/wirepair/ewserver/types"

// Storer interface for storing data
type Storer interface {
	Open(config *Config) error // Opens or creates the data store
	Close() error              // Closes the data store
	// User support
	StoreUser(user *types.User) error            // Stores a User
	DeleteUserByID(ID []byte) error              // Deletes a User
	FindUserByID(ID []byte) (*types.User, error) // Finds a User by ID
	// APIUser support
	StoreAPIUser(apiUser *types.APIUser) error // Stores an APIUser
	DeleteAPIUserByID(ID []byte) error         // Deletes an APIUser by ID
}
