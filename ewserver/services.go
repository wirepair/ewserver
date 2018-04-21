package ewserver

// Services is a simple container of our various domain services
type Services struct {
	UserService    UserService
	APIUserService APIUserService
	RoleService    RoleService
	LogService     LogService
}

// NewServices adds the various services to the Services container
func NewServices(userService UserService, apiUserService APIUserService, roleService RoleService, logService LogService) *Services {
	return &Services{UserService: userService, APIUserService: apiUserService, RoleService: roleService, LogService: logService}
}
