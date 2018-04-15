package mock

import "net/http"

// Sessions mocks our session implementation
type Sessions struct {
	AddFn      func(w http.ResponseWriter, req *http.Request, key string, value interface{}) error
	AddInvoked bool

	DestroyFn      func(w http.ResponseWriter, req *http.Request) error
	DestroyInvoked bool

	RenewFn      func(w http.ResponseWriter, req *http.Request) error
	RenewInvoked bool

	GetStringFn      func(req *http.Request, key string) string
	GetStringInvoked bool

	PopStringFn      func(w http.ResponseWriter, req *http.Request, key string) string
	PopStringInvoked bool

	LoadFn      func(req *http.Request, key string, result interface{}) error
	LoadInvoked bool

	PopLoadFn      func(w http.ResponseWriter, req *http.Request, key string, result interface{}) error
	PopLoadInvoked bool
}

// Destroy the session
func (s *Sessions) Destroy(w http.ResponseWriter, req *http.Request) error {
	s.DestroyInvoked = true
	return s.DestroyFn(w, req)
}

// Renew the session token
func (s *Sessions) Renew(w http.ResponseWriter, req *http.Request) error {
	s.RenewInvoked = true
	return s.RenewFn(w, req)
}

// Add a value to this session
func (s *Sessions) Add(w http.ResponseWriter, req *http.Request, key string, value interface{}) error {
	s.AddInvoked = true
	return s.AddFn(w, req, key, value)
}

// GetString value from this session
func (s *Sessions) GetString(req *http.Request, key string) string {
	s.GetStringInvoked = true
	return s.GetString(req, key)
}

// PopString pops a string value from our session, removing it and returning to caller
func (s *Sessions) PopString(w http.ResponseWriter, req *http.Request, key string) string {
	s.PopStringInvoked = true
	return s.PopStringFn(w, req, key)
}

// Load a value from the session into the result interface
func (s *Sessions) Load(req *http.Request, key string, result interface{}) error {
	s.LoadInvoked = true
	return s.LoadFn(req, key, result)
}

// PopLoad pops a value into the result interface and removes it from the session
func (s *Sessions) PopLoad(w http.ResponseWriter, req *http.Request, key string, result interface{}) error {
	s.PopLoadInvoked = true
	return s.PopLoadFn(w, req, key, result)
}
