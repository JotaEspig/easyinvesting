package login

import (
	"easyinvesting/config"
	"easyinvesting/pkg/controller"
	"easyinvesting/pkg/repository"
	"easyinvesting/pkg/service"
	"easyinvesting/pkg/types"
)

var AvailableRoutes []types.Route
var userController *controller.UserController

func init() {
	repo := repository.NewUserRepository(config.DB())
	service := service.NewUserService(repo)
	userController = controller.NewUserController(service)

	AvailableRoutes = []types.Route{
		{Method: "POST", Path: "/login", Fn: userController.Login()},
		{Method: "POST", Path: "/signup", Fn: userController.Register()},
	}
}
