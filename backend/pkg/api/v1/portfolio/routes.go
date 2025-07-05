package portfolio

import (
	"easyinvesting/config"
	"easyinvesting/pkg/controller"
	"easyinvesting/pkg/repository"
	"easyinvesting/pkg/service"
	"easyinvesting/pkg/types"

	"github.com/labstack/echo/v4"
)

var AvailableRoutes []types.Route
var assetController *controller.AssetController

func init() {
	repo := repository.NewAssetRepository(config.DB())
	service := service.NewAssetService(repo)
	assetController = controller.NewAssetController(service)

	AvailableRoutes = []types.Route{
		{Method: "POST", Path: "/asset/add", Fn: assetController.AddUserAsset(), Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
		{Method: "GET", Path: "/asset/list", Fn: assetController.GetPaginatedUserAssets(), Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
		{Method: "GET", Path: "/asset/:id", Fn: assetController.GetUserAsset(), Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
		{Method: "GET", Path: "/asset/:id/realtime", Fn: GetRealTimeAssetData, Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
		{Method: "POST", Path: "/realtimeupdate", Fn: UpdateRealTimeAssetsData, Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
		{Method: "POST", Path: "/asset/entry/add", Fn: AddUserAssetEntry, Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
		{Method: "GET", Path: "/asset/entry/:id", Fn: GetUserAssetEntry, Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
	}
}
