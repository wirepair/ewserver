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
	adminRoutes.GET("/all", AdminUsers(userService, e))
	adminRoutes.GET("/:user/details", AdminUser(userService, e))
	adminRoutes.PUT("/:user/create", AdminCreateUser(userService, e))
	adminRoutes.DELETE("/:user/delete", AdminDeleteUser(userService, e))
}

// RegisterAdminAPIRoutes for managing API users
func RegisterAdminAPIRoutes(apiUserService ewserver.APIUserService, routes *gin.RouterGroup, e *gin.Engine) {
	adminRoutes := routes.Group("/admin/api_users")
	adminRoutes.GET("/:key/details", AdminAPIUser(apiUserService, e))
	adminRoutes.GET("/all", AdminAPIUser(apiUserService, e))
	adminRoutes.PUT("/add", AdminAddAPIUser(apiUserService, e))
	adminRoutes.DELETE("/api_user/:id/delete", AdminDeleteAPIUser(apiUserService, e))
}
