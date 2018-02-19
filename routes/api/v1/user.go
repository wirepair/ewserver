package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/store"
)

// AddUser to allow access to the UI
func AddUser(dataStore store.Storer, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, nil)
	}
}

// UpdateUser for UI access
func UpdateUser(dataStore store.Storer, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, nil)
	}
}

// DeleteUser from UI access
func DeleteUser(dataStore store.Storer, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, nil)
	}
}
