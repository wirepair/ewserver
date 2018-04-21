package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/api/v1"
	"github.com/wirepair/ewserver/internal/authz"
)

// Require authorization token (api or session)
func Require(authorizer authz.Authorizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authorizer.Authorize(c.Request) {
			c.Next()
			return
		}

		c.Redirect(301, v1.LoginPath)
		c.AbortWithStatus(301)
	}
}
