package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/api/v1"
	"github.com/wirepair/ewserver/internal/authz"
)

// Require authorization token (api or session) and the specified role
func Require(authorizer authz.Authorizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("authorizing request...\n")
		if authorizer.Authorize(c.Request) {
			c.Next()
			return
		}
		log.Printf("not authorized, redirecting...\n")
		c.Redirect(301, v1.LoginPath)
		c.AbortWithStatus(301)
	}
}
