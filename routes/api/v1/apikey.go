package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/store"
)

// GenerateKey for a User
func GenerateKey(dataStore store.Storer, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, nil)
	}
}

// RevokeKey for a User
func RevokeKey(dataStore store.Storer, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, nil)
	}
}
