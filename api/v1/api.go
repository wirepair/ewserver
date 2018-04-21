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

// RegisterAdminRoutes for managing the system
func RegisterAdminRoutes(services *ewserver.Services, e *gin.Engine) {
	// setup admin routes
	apiRoutes := e.Group("v1")
	userRoutes := apiRoutes.Group("/admin/users")
	userRoutes.GET("/details/:user", AdminUserDetails(services.UserService, services.LogService, e))
	userRoutes.GET("/list", AdminUsersDetails(services.UserService, services.LogService, e))
	userRoutes.PUT("/create", AdminCreateUser(services.UserService, services.LogService, e))
	userRoutes.POST("/reset_password", AdminResetPassword(services.UserService, services.LogService, e))
	userRoutes.DELETE("/delete/:user", AdminDeleteUser(services.UserService, services.LogService, e))

	apiAdminRoutes := apiRoutes.Group("/admin/api_users")
	apiAdminRoutes.GET("/details/:id", AdminAPIUserDetails(services.APIUserService, services.LogService, e))
	apiAdminRoutes.GET("/list", AdminAPIUsersDetails(services.APIUserService, services.LogService, e))
	apiAdminRoutes.PUT("/create", AdminCreateAPIUser(services.APIUserService, services.LogService, e))
	apiAdminRoutes.DELETE("/delete/:id", AdminDeleteAPIUser(services.APIUserService, services.LogService, e))

	roleRoutes := apiRoutes.Group("/admin/roles")
	roleRoutes.GET("/list", AdminRoleList(services.RoleService, services.LogService, e))
	roleRoutes.PUT("/permissions", AdminAddPermission(services.RoleService, services.LogService, e))
	roleRoutes.DELETE("/permissions", AdminDeletePermission(services.RoleService, services.LogService, e))
	roleRoutes.POST("/group", AdminAddUserToGroup(services.RoleService, services.LogService, e))
	roleRoutes.DELETE("/group", AdminDeleteUserFromGroup(services.RoleService, services.LogService, e))
}

// RegisterAuthnRoutes registers the authentication (login/logout) routes under /user
func RegisterAuthnRoutes(authnService ewserver.AuthnService, logService ewserver.LogService, e *gin.Engine) {
	routes := e.Group("/")
	e.LoadHTMLGlob("../../web/templates/**/*")
	routes.GET(LoginPath, LoginPage(e))
	routes.POST(LoginPath, Authenticate(authnService, logService, e))
	routes.GET("/logout", Logout(authnService, logService, e))
}
