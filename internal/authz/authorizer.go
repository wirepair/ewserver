package authz

import "net/http"

// Authorizer handles authorizing HTTP requests
type Authorizer interface {
	Authorize(r *http.Request) bool
	APIAuthorize(r *http.Request, apiKey string) bool    // APIAuthorize for an API User
	UserAuthorize(r *http.Request, username string) bool // UserAuthorize for regular users
}
