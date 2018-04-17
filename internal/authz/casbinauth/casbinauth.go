package casbinauth

import (
	"net/http"

	"github.com/casbin/casbin"
	"github.com/wirepair/ewserver/ewserver"
	"github.com/wirepair/ewserver/internal/session"
)

// CasbinAuthorizer uses casbin to authorize requests
type CasbinAuthorizer struct {
	enforcer       *casbin.SyncedEnforcer
	apiUserService ewserver.APIUserService
	sessions       session.Manager
}

// New returns a new CasbinAuthorizer
func NewAuthorizer(enforcer *casbin.SyncedEnforcer, apiUserService ewserver.APIUserService, sessions session.Manager) *CasbinAuthorizer {
	return &CasbinAuthorizer{enforcer: enforcer, apiUserService: apiUserService, sessions: sessions}
}

// Authorize validates the user data from a request is authorized to access a resource
func (a *CasbinAuthorizer) Authorize(r *http.Request) bool {
	// Check API based auth first
	apiKey := r.Header.Get(ewserver.APIKeyHeader)

	// if apikey is not empty and they are authorized, return true.
	if apiKey != "" && a.APIAuthorize(r, apiKey) {
		return true
	}

	user := &ewserver.User{}
	if err := a.sessions.Load(r, "user", user); err != nil {
		return false
	}
	return a.UserAuthorize(r, string(user.UserName))
}

// APIAuthorize for API Users
func (a *CasbinAuthorizer) APIAuthorize(r *http.Request, apiKey string) bool {
	user, err := a.apiUserService.APIUser(ewserver.APIKey(apiKey))
	if err != nil || user.Name == "" {
		return false
	}

	subject := user.Name
	object := r.URL.Path
	action := r.Method
	return a.enforcer.Enforce(subject, object, action)
}

// UserAuthorize for regular users
func (a *CasbinAuthorizer) UserAuthorize(r *http.Request, username string) bool {
	subject := username
	object := r.URL.Path
	action := r.Method
	return a.enforcer.Enforce(subject, object, action)
}
