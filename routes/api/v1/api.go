package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/repositories"
	"github.com/wirepair/ewserver/store"
)

func default404(e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, nil)
	}
}

// RegisterRoutesV1 for v1 of the API.
func RegisterRoutesV1(dataStore store.Storer, e *gin.Engine) {
	userRepo := repositories.NewUserStore(dataStore)

	routes := e.Group("/v1")
	routes.GET("/", default404(e))

	adminRoutes := routes.Group("/admin")
	adminRoutes.GET("/user/details", AdminGetUserDetails(userRepo, e))
	adminRoutes.GET("/user/all/details", AdminGetAllUsersDetails(userRepo, e))
	adminRoutes.PUT("/user/add", AdminAddUser(userRepo, e))
	adminRoutes.POST("/user/edit", AdminEditUser(userRepo, e))
	adminRoutes.DELETE("/user/delete", AdminDeleteUser(userRepo, e))
}
