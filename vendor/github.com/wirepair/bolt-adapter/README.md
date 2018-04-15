Bolt Adapter [![Build Status](https://travis-ci.org/wirepair/bolt-adapter.svg?branch=master)](https://travis-ci.org/wirepair/bolt-adapter) [![Godoc](https://godoc.org/github.com/wirepair/bolt-adapter?status.svg)](https://godoc.org/github.com/wirepair/bolt-adapter)
====

Bolt Adapter is the [BoltDB](https://github.com/boltdb/bolt) adapter for [Casbin](https://github.com/casbin/casbin). With this library, Casbin can load a policy from a previously opened BoltDB or save the policy to it.

## Installation

    go get github.com/wirepair/bolt-adapter

## Simple Example

```go
package main

import (
	"flag"
	"log"

	"github.com/boltdb/bolt"
	"github.com/casbin/casbin"
	"github.com/casbin/casbin/persist"
	"github.com/wirepair/bolt-adapter"
)

var populate bool

func init() {
	flag.BoolVar(&populate, "populate", false, "populate the db from a file first")
}

func main() {
	flag.Parse()

	// Initialize a bolt DB adapter and use it in a Casbin enforcer:
	db, err := bolt.Open("db.dat", 0600, nil)
	if err != nil {
		log.Fatalf("error opening db: %s\n", err)
	}
	defer db.Close()

	adapter := boltadapter.NewAdapter(db) // Pass in the already open bolt DB.

	if populate {
		populateDB(adapter)
	}

	e := casbin.NewEnforcer("examples/rbac_model.conf", adapter)

	// Load the policy from DB.
	e.LoadPolicy()

	// Check the permission.
	e.Enforce("alice", "data1", "read")

	// Save the policy back to DB.
	e.SavePolicy()
}

// populateDB loads a policy from file and saves it to the DB
func populateDB(adapter persist.Adapter) {
	e := casbin.NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")
	adapter.SavePolicy(e.GetModel())
}
```

## Getting Help

- [Casbin](https://github.com/casbin/casbin)

## License

This project is under Apache 2.0 License. See the [LICENSE](LICENSE) file for the full license text.

## Author 
Isaac Dawson
