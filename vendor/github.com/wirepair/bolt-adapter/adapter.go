package boltadapter

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
)

var (
	// ErrEmptyPolicy will be returned if the bucket doesn't have any policy data
	ErrEmptyPolicy = errors.New("policy was empty")
)

// CasbinRule represents a policy type and their values
type CasbinRule struct {
	PType string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}

// Adapter represents the BoltDB adapter for policy storage.
type Adapter struct {
	db  *bolt.DB
	key []byte
}

func newAdapter(db *bolt.DB, key string) *Adapter {
	a := &Adapter{}
	a.db = db
	a.key = []byte(key)

	a.open()

	return a
}

// NewAdapter is the constructor for Adapter. Assumes the bolt db is already opened.
func NewAdapter(db *bolt.DB) *Adapter {
	return newAdapter(db, "casbin_rules")
}

// NewBoltAdapter is the constructor for Adapter. Assumes the bolt db is already opened.
func NewBoltAdapter(db *bolt.DB, key string) *Adapter {
	return newAdapter(db, key)
}

func (a *Adapter) open() {
	err := a.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(a.key))
		return err
	})

	// i don't like panic'ing here but that's what the other adapters do.
	if err != nil {
		panic(err)
	}
}

func loadPolicyLine(line CasbinRule, model model.Model) {
	lineText := line.PType
	if line.V0 != "" {
		lineText += ", " + line.V0
	}
	if line.V1 != "" {
		lineText += ", " + line.V1
	}
	if line.V2 != "" {
		lineText += ", " + line.V2
	}
	if line.V3 != "" {
		lineText += ", " + line.V3
	}
	if line.V4 != "" {
		lineText += ", " + line.V4
	}
	if line.V5 != "" {
		lineText += ", " + line.V5
	}

	persist.LoadPolicyLine(lineText, model)
}

// LoadPolicy loads policy from database.
func (a *Adapter) LoadPolicy(model model.Model) error {
	return a.db.View(func(tx *bolt.Tx) error {
		lines := make([]CasbinRule, 0)
		bucket := tx.Bucket([]byte(a.key))
		policy := bucket.Get([]byte("policy"))
		if policy == nil {
			return ErrEmptyPolicy
		}

		if err := json.Unmarshal(policy, &lines); err != nil {
			return err
		}

		for _, line := range lines {
			loadPolicyLine(line, model)
		}

		return nil
	})
}

func savePolicyLine(ptype string, rule []string) CasbinRule {
	line := CasbinRule{}

	line.PType = ptype
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}

// SavePolicy saves policy to database.
func (a *Adapter) SavePolicy(model model.Model) error {

	var lines []CasbinRule

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, line)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, line)
		}
	}

	text, err := json.Marshal(lines)
	if err != nil {
		return err
	}

	return a.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(a.key))
		return bucket.Put([]byte("policy"), []byte(text))
	})
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}
