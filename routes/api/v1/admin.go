package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/repositories"
)

// AdminAddUser to allow access to the UI
func AdminAddUser(userRepo repositories.UserRepositorer, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, nil)
	}
}

// AdminGetAllUsersDetails returns details for all users
func AdminGetAllUsersDetails(userRepo repositories.UserRepositorer, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, nil)
	}
}

// AdminGetUserDetails returns details for a specified user
func AdminGetUserDetails(userRepo repositories.UserRepositorer, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, nil)
	}
}

// AdminEditUser for UI access
func AdminEditUser(userRepo repositories.UserRepositorer, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, nil)
	}
}

// AdminDeleteUser from UI access
func AdminDeleteUser(userRepo repositories.UserRepositorer, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, nil)
	}
}
