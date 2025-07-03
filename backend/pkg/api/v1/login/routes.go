package login

import (
	"easyinvesting/config"
	"easyinvesting/pkg/controllers"
	"easyinvesting/pkg/repositories"
	"easyinvesting/pkg/services"
	"easyinvesting/pkg/types"
)

var AvailableRoutes []types.Route
var userController *controllers.UserController

func init() {
	repo := repositories.NewUserRepository(config.DB())
	service := services.NewUserService(repo)
	userController = controllers.NewUserController(service)

	AvailableRoutes = []types.Route{
		{Method: "POST", Path: "/login", Fn: userController.Login()},
		{Method: "POST", Path: "/signup", Fn: userController.Register()},
	}
}
