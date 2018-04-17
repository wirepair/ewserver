package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/ewserver"
)

// AdminGroups ...
func AdminGroups(roleService ewserver.RoleService, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		groups := roleService.Groups()

		c.JSON(200, gin.H{"groups": groups})
	}
}
