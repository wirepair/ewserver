package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/ewserver"
	"github.com/wirepair/ewserver/internal/converter"
)

// AdminRoleList users and groups
func AdminRoleList(roleService ewserver.RoleService, logService ewserver.LogService, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleNames := roleService.RoleNames()
		roleMap := roleService.RoleMap()
		permissions := roleService.Permissions()
		c.JSON(200, gin.H{"status": "OK", "role_names": roleNames, "role_map": roleMap, "permissions": permissions})
	}
}

// AdminAddPermission add a new permission to user or group
func AdminAddPermission(roleService ewserver.RoleService, logService ewserver.LogService, e *gin.Engine) gin.HandlerFunc {
	type permission struct {
		Subject  string   `json:"subject"`
		Resource string   `json:"resource"`
		Action   []string `json:"actions"`
	}

	return func(c *gin.Context) {
		perm := &permission{}
		if err := c.BindJSON(perm); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		regex := converter.ActionsRegex(perm.Action)
		if regex == "" {
			c.JSON(500, gin.H{"error": "invalid action specified"})
			return
		}
		err := roleService.AddPermission(perm.Subject, perm.Resource, regex)
		defaultReturn(err, c)
	}
}

// AdminDeletePermission delete a permission
func AdminDeletePermission(roleService ewserver.RoleService, logService ewserver.LogService, e *gin.Engine) gin.HandlerFunc {
	type permission struct {
		Subject  string `json:"subject"`
		Resource string `json:"resource"`
		Action   string `json:"actions"`
	}

	return func(c *gin.Context) {
		perm := &permission{}
		if err := c.BindJSON(perm); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		err := roleService.DeletePermission(perm.Subject, perm.Resource, perm.Action)
		logService.Info("perm delete", err.Error())
		defaultReturn(err, c)
	}
}

// AdminAddSubjectToRole add a subject to a role (such as user to group)
func AdminAddSubjectToRole(roleService ewserver.RoleService, logService ewserver.LogService, e *gin.Engine) gin.HandlerFunc {
	type role struct {
		User string `json:"user"`
		Role string `json:"role"`
	}

	return func(c *gin.Context) {
		addRole := &role{}
		if err := c.BindJSON(addRole); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		err := roleService.AddSubjectToRole(addRole.User, addRole.Role)
		defaultReturn(err, c)
	}
}

// AdminDeleteSubjectFromRole delete a user from a group
func AdminDeleteSubjectFromRole(roleService ewserver.RoleService, logService ewserver.LogService, e *gin.Engine) gin.HandlerFunc {
	type role struct {
		User string `json:"user"`
		Role string `json:"role"`
	}

	return func(c *gin.Context) {
		deleteRole := &role{}
		if err := c.BindJSON(deleteRole); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		err := roleService.DeleteSubjectFromRole(deleteRole.User, deleteRole.Role)
		defaultReturn(err, c)
	}
}
