package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/ewserver"
)

// AdminRoleList users and groups
func AdminRoleList(roleService ewserver.RoleService, logService ewserver.LogService, e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		groups := roleService.Groups()
		users := roleService.Users()
		permissions := roleService.Permissions()
		c.JSON(200, gin.H{"groups": groups, "users": users, "permissions": permissions})
	}
}

// AdminAddPermission add a new permission to user or group
func AdminAddPermission(roleService ewserver.RoleService, logService ewserver.LogService, e *gin.Engine) gin.HandlerFunc {
	type permission struct {
		Subject string `json:"subject"`
		Object  string `json:"object"`
		Method  string `json:"method"`
	}

	return func(c *gin.Context) {
		perm := &permission{}
		if err := c.BindJSON(perm); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		err := roleService.AddPermission(perm.Subject, perm.Object, perm.Method)
		defaultReturn(err, c)
	}
}

// AdminDeletePermission delete a permission
func AdminDeletePermission(roleService ewserver.RoleService, logService ewserver.LogService, e *gin.Engine) gin.HandlerFunc {
	type permission struct {
		Subject string `json:"subject"`
		Object  string `json:"object"`
		Method  string `json:"method"`
	}

	return func(c *gin.Context) {
		perm := &permission{}
		if err := c.BindJSON(perm); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		err := roleService.DeletePermission(perm.Subject, perm.Object, perm.Method)
		defaultReturn(err, c)

	}
}

// AdminAddUserToGroup add a user to a group
func AdminAddUserToGroup(roleService ewserver.RoleService, logService ewserver.LogService, e *gin.Engine) gin.HandlerFunc {
	type role struct {
		User  string `json:"user"`
		Group string `json:"group"`
	}

	return func(c *gin.Context) {
		addRole := &role{}
		if err := c.BindJSON(addRole); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		err := roleService.AddUserToGroup(addRole.User, addRole.Group)
		defaultReturn(err, c)
	}
}

// AdminDeleteUserFromGroup delete a user from a group
func AdminDeleteUserFromGroup(roleService ewserver.RoleService, logService ewserver.LogService, e *gin.Engine) gin.HandlerFunc {
	type role struct {
		User  string `json:"user"`
		Group string `json:"group"`
	}

	return func(c *gin.Context) {
		deleteRole := &role{}
		if err := c.BindJSON(deleteRole); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		err := roleService.DeleteUserFromGroup(deleteRole.User, deleteRole.Group)
		defaultReturn(err, c)
	}
}
