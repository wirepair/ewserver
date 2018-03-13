package v1

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/ewserver"
)

// AdminUserDetails returns the details of a single user
func AdminUserDetails(userService ewserver.UserService, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		userName := c.Param("user")
		user, err := userService.User(ewserver.UserName(userName))
		if err != nil {
			c.JSON(500, gin.H{"error": err})
		}

		// Just incase
		user.Password = []byte{}

		c.JSON(200, gin.H{"user": user})
	}
}

// AdminUsersDetails returns the details of all users
func AdminUsersDetails(userService ewserver.UserService, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := userService.Users()
		if err != nil {
			c.JSON(500, gin.H{"error": err})
		}

		for _, user := range users {
			// Just incase
			user.Password = []byte{}
		}

		c.JSON(200, gin.H{"users": users})
	}
}

// AdminCreateUser creates a new user.
func AdminCreateUser(userService ewserver.UserService, e *gin.Engine) gin.HandlerFunc {
	// embed the ewserver.User but expose/override the Password as a string to allow it to be read
	// from JSON
	type newUser struct {
		ewserver.User
		Password string `json:"password"`
	}

	return func(c *gin.Context) {
		user := &newUser{}
		if err := c.BindJSON(user); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		log.Printf("%#v\n", user)

		if err := userService.Create(&user.User, user.Password); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		c.JSON(200, gin.H{"status": "OK"})
	}
}

// AdminResetPassword changes the password for the specified user
func AdminResetPassword(userService ewserver.UserService, e *gin.Engine) gin.HandlerFunc {
	type passwordReset struct {
		UserName    ewserver.UserName `json:"user_name"`
		NewPassword string            `json:"password"`
	}

	return func(c *gin.Context) {
		passwordRequest := &passwordReset{}

		if err := c.BindJSON(passwordRequest); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		if err := userService.ResetPassword(passwordRequest.UserName, passwordRequest.NewPassword); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}
		c.JSON(200, gin.H{"status": "OK"})
	}
}

// AdminDeleteUser deletes a user
func AdminDeleteUser(userService ewserver.UserService, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		userName := c.Param("user")
		if err := userService.Delete(ewserver.UserName(userName)); err != nil {
			c.JSON(500, gin.H{"error": err})
		}
		c.JSON(200, gin.H{"status": "OK"})
	}
}

// AdminAPIUserDetails returns the details of an API User
func AdminAPIUserDetails(apiUserService ewserver.APIUserService, e *gin.Engine) gin.HandlerFunc {

	return func(c *gin.Context) {
		id := c.Param("id")
		user, err := apiUserService.APIUserByID([]byte(id))
		if err != nil {
			c.JSON(500, gin.H{"error": err})
		}

		c.JSON(200, gin.H{"user": user})
	}
}

// AdminAPIUsersDetails returns the details of all API Users
func AdminAPIUsersDetails(apiUserService ewserver.APIUserService, e *gin.Engine) gin.HandlerFunc {

	return func(c *gin.Context) {
		apiUsers, err := apiUserService.APIUsers()
		if err != nil {
			c.JSON(500, gin.H{"error": err})
		}

		c.JSON(200, gin.H{"users": apiUsers})
	}
}

// AdminCreateAPIUser adds a new API User, generates a new key prior to creating.
func AdminCreateAPIUser(apiUserService ewserver.APIUserService, e *gin.Engine) gin.HandlerFunc {

	return func(c *gin.Context) {
		apiUser := &ewserver.APIUser{}
		if err := c.BindJSON(apiUser); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		key, err := ewserver.GenerateAPIKey()
		if err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		apiUser.Key = key

		if err := apiUserService.Create(apiUser); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}
		c.JSON(200, gin.H{"status": "OK"})
	}
}

// AdminDeleteAPIUser deletes the API key by first looking up the ID to get the APIKey.
func AdminDeleteAPIUser(apiUserService ewserver.APIUserService, e *gin.Engine) gin.HandlerFunc {

	return func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("finding %s with %#v\n", id, []byte(id))
		apiUser, err := apiUserService.APIUserByID([]byte(id))
		if err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		if err := apiUserService.Delete(apiUser.Key); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}
		c.JSON(200, gin.H{"status": "OK"})
	}
}
