package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/ewserver"
)

func default404(e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, nil)
	}
}

// RegisterAdminRoutes for managing Users
func RegisterAdminRoutes(userService ewserver.UserService, routes *gin.RouterGroup, e *gin.Engine) {
	adminRoutes := routes.Group("/admin/users")
	adminRoutes.GET("/details/:user", AdminUserDetails(userService, e))
	adminRoutes.GET("/all_details", AdminUsersDetails(userService, e))
	adminRoutes.PUT("/create", AdminCreateUser(userService, e))
	adminRoutes.POST("/reset_password", AdminResetPassword(userService, e))
	adminRoutes.DELETE("/delete/:user", AdminDeleteUser(userService, e))
}

// RegisterAdminAPIRoutes for managing APIUsers
func RegisterAdminAPIRoutes(apiUserService ewserver.APIUserService, routes *gin.RouterGroup, e *gin.Engine) {
	adminRoutes := routes.Group("/admin/api_users")
	adminRoutes.GET("/details/:id", AdminAPIUserDetails(apiUserService, e))
	adminRoutes.GET("/all_details", AdminAPIUsersDetails(apiUserService, e))
	adminRoutes.PUT("/create", AdminCreateAPIUser(apiUserService, e))
	adminRoutes.DELETE("/delete/:id", AdminDeleteAPIUser(apiUserService, e))
}

func RegisterAuthnRoutes(authnService ewserver.AuthnService, routes *gin.RouterGroup, e *gin.Engine) {
	authRoutes := routes.Group("/user")
	authRoutes.POST("/login", Auth(authnService, e))
	authRoutes.GET("/logout", Logout(authnService, e))
}
