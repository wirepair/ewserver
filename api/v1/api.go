package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/ewserver"
)

const (
	// LoginPath for login GET/POST
	LoginPath = "/login"
)

func default404(e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, nil)
	}
}

// RegisterAdminRoutes for managing Users under /admin/users
func RegisterAdminRoutes(userService ewserver.UserService, routes *gin.RouterGroup, e *gin.Engine) {
	routes.GET("/details/:user", AdminUserDetails(userService, e))
	routes.GET("/all_details", AdminUsersDetails(userService, e))
	routes.PUT("/create", AdminCreateUser(userService, e))
	routes.POST("/reset_password", AdminResetPassword(userService, e))
	routes.DELETE("/delete/:user", AdminDeleteUser(userService, e))
}

// RegisterAdminAPIRoutes for managing APIUsers under admin/api_users
func RegisterAdminAPIRoutes(apiUserService ewserver.APIUserService, routes *gin.RouterGroup, e *gin.Engine) {
	routes.GET("/details/:id", AdminAPIUserDetails(apiUserService, e))
	routes.GET("/all_details", AdminAPIUsersDetails(apiUserService, e))
	routes.PUT("/create", AdminCreateAPIUser(apiUserService, e))
	routes.DELETE("/delete/:id", AdminDeleteAPIUser(apiUserService, e))
}

// RegisterAuthnRoutes registers the authentication (login/logout) routes under /user
func RegisterAuthnRoutes(authnService ewserver.AuthnService, routes *gin.RouterGroup, e *gin.Engine) {
	routes.GET(LoginPath, LoginPage(e))
	routes.POST(LoginPath, Authenticate(authnService, e))
	routes.GET("/logout", Logout(authnService, e))
}
