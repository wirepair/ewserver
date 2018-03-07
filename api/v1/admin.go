package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/ewserver"
)

// AdminChangePassword changes the password for the specified user
func AdminChangePassword(userService ewserver.UserService, e *gin.Engine) gin.HandlerFunc {
	//const op errors.Op = "routes/AdminAddUser"
	return func(c *gin.Context) {
		var newUser ewserver.User
		if err := c.BindJSON(newUser); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		/*
			if err := userRepo.AddUser(&newUser); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		*/

		c.JSON(200, gin.H{"status": "OK"})
	}
}

// AdminAddAPIUser adds a new API Key
func AdminAddAPIUser(userService ewserver.UserService, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser ewserver.User
		if err := c.BindJSON(newUser); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		/*
			if err := userRepo.AddUser(&newUser); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		*/

		c.JSON(200, gin.H{"status": "OK"})
	}
}

// AdminDeleteAPIUser deletes the API key
func AdminDeleteAPIUser(userService ewserver.UserService, e *gin.Engine) gin.HandlerFunc {
	//const op errors.Op = "routes/AdminAddUser"
	return func(c *gin.Context) {
		var apiUser ewserver.APIUser
		if err := c.BindJSON(apiUser); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		/*
			if err := userRepo.AddUser(&newUser); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		*/

		c.JSON(200, gin.H{"status": "OK"})
	}
}
