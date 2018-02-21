package store

import "github.com/wirepair/ewserver/types"

// Storer interface for storing data
type Storer interface {
	Open(config *Config) error                                       // Opens or creates the data store
	Close() error                                                    // Closes the data store
	StoreUser(user *types.User) error                                // Stores a User
	DeleteUserByName(userName types.UserName) error                  // Deletes a User by UserName
	FindUserByUserName(userName types.UserName) (*types.User, error) // Finds a User by UserName
	FindUserByID(ID []byte) (*types.User, error)                     // Finds a User by ID
	FindUserByAPIKey(Key []byte) (*types.User, error)                // Finds a User by the API Key
	FindAllUsers() ([]*types.User, error)                            // Finds and returns all Users
}
