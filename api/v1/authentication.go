package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/ewserver"
)

// Auth a user to create a session
func Auth(authnService ewserver.AuthnService, e *gin.Engine) gin.HandlerFunc {
	type login struct {
		UserName ewserver.UserName
		Password string
	}

	return func(c *gin.Context) {
		attempt := &login{}
		if err := c.BindJSON(attempt); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		user, err := authnService.Authenticate(attempt.UserName, attempt.Password)
		if err != nil {
			c.JSON(500, gin.H{"error": err})
		}

		c.JSON(200, gin.H{"user": user})
	}
}

// Logout a user by destroying their session
func Logout(authnService ewserver.AuthnService, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	}
}
