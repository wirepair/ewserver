package scs

import (
	"net/http"

	"github.com/alexedwards/scs"
)

// Sessions using scs
type Sessions struct {
	*scs.Manager
}

// New creates a new session manager backed by scs
func New(manager *scs.Manager) *Sessions {
	return &Sessions{Manager: manager}
}

// Destroy the session
func (s Sessions) Destroy(w http.ResponseWriter, req *http.Request) error {
	scssession := s.Manager.Load(req)
	return scssession.Destroy(w)
}

// Renew the session token
func (s Sessions) Renew(w http.ResponseWriter, req *http.Request) error {
	scssession := s.Manager.Load(req)
	return scssession.RenewToken(w)
}

// Add a value to this session
func (s Sessions) Add(w http.ResponseWriter, req *http.Request, key string, value interface{}) error {
	scssession := s.Manager.Load(req)
	if str, ok := value.(string); ok {
		return scssession.PutString(w, key, str)
	}
	return scssession.PutObject(w, key, value)
}

// GetString value from this session
func (s Sessions) GetString(w http.ResponseWriter, req *http.Request, key string) string {
	scssession := s.Manager.Load(req)
	result, _ := scssession.GetString(key)
	return result
}

// PopString pops a string value from our session, removing it and returning to caller
func (s Sessions) PopString(w http.ResponseWriter, req *http.Request, key string) string {
	scssession := s.Manager.Load(req)
	result, _ := scssession.PopString(w, key)
	return result
}

// Load a value from the session into the result interface
func (s Sessions) Load(req *http.Request, key string, result interface{}) error {
	scssession := s.Manager.Load(req)
	return scssession.GetObject(key, result)
}

// PopLoad pops a value into the result interface and removes it from the session
func (s Sessions) PopLoad(w http.ResponseWriter, req *http.Request, key string, result interface{}) error {
	scssession := s.Manager.Load(req)
	return scssession.PopObject(w, key, result)
}
