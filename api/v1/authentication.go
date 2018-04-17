package v1

import (
	"net/http"

	"github.com/wirepair/ewserver/internal/session"

	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/ewserver"
)

// LoginPage displays the login page to the user
func LoginPage(e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "user/login.tmpl", gin.H{
			"title": "Login",
		})
	}
}

// Authenticate a user to create a session
func Authenticate(authnService ewserver.AuthnService, e *gin.Engine) gin.HandlerFunc {
	type login struct {
		UserName ewserver.UserName
		Password string
	}

	return func(c *gin.Context) {
		sessions := c.MustGet("sessions").(session.Manager)
		attempt := &login{}

		if err := c.BindJSON(attempt); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		user, err := authnService.Authenticate(attempt.UserName, attempt.Password)
		if err != nil {
			c.JSON(401, gin.H{"error": err})
			return
		}

		sessions.Renew(c.Writer, c.Request)
		sessions.Add(c.Writer, c.Request, "user", user)
		c.JSON(200, gin.H{"user": user})
	}
}

// Logout a user by destroying their session
func Logout(authnService ewserver.AuthnService, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessions := c.MustGet("sessions").(session.Manager)
		sessions.Destroy(c.Writer, c.Request)
		c.JSON(200, gin.H{"status": "OK"})
	}
}
