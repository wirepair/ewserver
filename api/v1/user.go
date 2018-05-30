package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/ewserver"
	"github.com/wirepair/ewserver/internal/session"
)

// UserProfile returns the user profile to caller
func UserProfile(userService ewserver.UserService, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessions := c.MustGet("sessions").(session.Manager)
		user := &ewserver.User{}
		if err := sessions.Load(c.Request, "user", user); err != nil {
			defaultReturn(err, c)
			return
		}

		c.JSON(200, gin.H{"status": "OK", "user": user})
	}
}
