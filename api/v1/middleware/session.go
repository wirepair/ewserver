package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/ewserver"
	"github.com/wirepair/ewserver/internal/session"
)

// EnsureSession verfies a session exists and that it's bound to a user
// (provided x-api-key does not exist)
// Even though we add the anonymous user to the session, it will not exist for
// the authorization check, so the first request will redirect to /login
// after issuing a new cookie.
func EnsureSession(sessions session.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request has API header first
		apiKey := c.GetHeader(ewserver.APIKeyHeader)
		if apiKey != "" {
			c.Next()
			return
		}

		user := &ewserver.User{}
		if err := sessions.Load(c.Request, "user", user); err != nil || user.UserName == "" {
			user.UserName = "anonymous"
			sessions.Add(c.Writer, c.Request, "user", user)
		}

		// Add the session to the context
		c.Set("sessions", sessions)
		c.Next()
	}
}
