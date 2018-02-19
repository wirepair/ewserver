package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/store"
)

func default404(e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, nil)
	}
}

// RegisterRoutes for the v1 API.
func RegisterRoutes(dataStore *store.Storer, e *gin.Engine) {
	routes := e.Group("/v1")
	routes.GET("/", default404(e))
}
