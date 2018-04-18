package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/ewserver"
)

const (
	// LoginPath for login GET/POST
	LoginPath = "/login"
)

func defaultReturn(err error, c *gin.Context) {
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{"status": "OK"})
}

// RegisterAdminRoutes for managing Users under /admin/users
func RegisterAdminRoutes(services *ewserver.Services, e *gin.Engine) {
	// setup admin routes
	apiRoutes := e.Group("v1")
	userRoutes := apiRoutes.Group("/admin/users")
	userRoutes.GET("/details/:user", AdminUserDetails(services.UserService, e))
	userRoutes.GET("/list", AdminUsersDetails(services.UserService, e))
	userRoutes.PUT("/create", AdminCreateUser(services.UserService, e))
	userRoutes.POST("/reset_password", AdminResetPassword(services.UserService, e))
	userRoutes.DELETE("/delete/:user", AdminDeleteUser(services.UserService, e))

	apiAdminRoutes := apiRoutes.Group("/admin/api_users")
	apiAdminRoutes.GET("/details/:id", AdminAPIUserDetails(services.APIUserService, e))
	apiAdminRoutes.GET("/list", AdminAPIUsersDetails(services.APIUserService, e))
	apiAdminRoutes.PUT("/create", AdminCreateAPIUser(services.APIUserService, e))
	apiAdminRoutes.DELETE("/delete/:id", AdminDeleteAPIUser(services.APIUserService, e))

	roleRoutes := apiRoutes.Group("/admin/roles")
	roleRoutes.GET("/list", AdminRoleList(services.RoleService, e))
	roleRoutes.PUT("/permissions", AdminAddPermission(services.RoleService, e))
	roleRoutes.DELETE("/permissions", AdminDeletePermission(services.RoleService, e))
	roleRoutes.POST("/group", AdminAddUserToGroup(services.RoleService, e))
	roleRoutes.DELETE("/group", AdminDeleteUserFromGroup(services.RoleService, e))

}

// RegisterAuthnRoutes registers the authentication (login/logout) routes under /user
func RegisterAuthnRoutes(authnService ewserver.AuthnService, e *gin.Engine) {
	routes := e.Group("/")
	e.LoadHTMLGlob("../../web/templates/**/*")
	routes.GET(LoginPath, LoginPage(e))
	routes.POST(LoginPath, Authenticate(authnService, e))
	routes.GET("/logout", Logout(authnService, e))
}
