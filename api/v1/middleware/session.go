package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/ewserver"
	"github.com/wirepair/ewserver/internal/session"
)

// EnsureSession exists and it's bound to user (provided x-api-key does not exist)
func EnsureSession(sessions session.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request has API header first
		apiKey := c.GetHeader(ewserver.APIKeyHeader)
		if apiKey != "" {
			c.Next()
			return
		}

		user := &ewserver.User{}
		if err := sessions.Load(c.Request, "user", user); err != nil {
			user.UserName = "anonymous"
			sessions.Add(c.Writer, c.Request, "user", user)
		}
		// Add the session to the context
		c.Set("sessions", sessions)
		c.Next()
	}
}
