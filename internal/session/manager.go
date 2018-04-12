package session

import "net/http"

// Manager provides session management features.
type Manager interface {
	Add(w http.ResponseWriter, req *http.Request, key string, value interface{}) error
	Destroy(w http.ResponseWriter, req *http.Request) error
	Renew(w http.ResponseWriter, req *http.Request) error
	GetString(req *http.Request, key string) string
	PopString(w http.ResponseWriter, req *http.Request, key string) string
	Load(req *http.Request, key string, result interface{}) error
	PopLoad(w http.ResponseWriter, req *http.Request, key string, result interface{}) error
}
