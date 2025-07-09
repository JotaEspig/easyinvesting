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
var assetEntryController *controller.AssetEntryController
var assetOnMarketController *controller.AssetOnMarketController

func init() {
	assetRepo := repository.NewAssetRepository(config.DB())
	assetService := service.NewAssetService(assetRepo)
	entryRepo := repository.NewAssetEntryRepository(config.DB())
	entryService := service.NewAssetEntryService(assetRepo, entryRepo)
	assetOnMarketRepo := repository.NewAssetOnMarketRepository(config.DB())
	dailyAssetRepo := repository.NewDailyAssetPriceRepository(config.DB())
	assetOnMarketService := service.NewAssetOnMarketService(assetOnMarketRepo, dailyAssetRepo)

	assetController = controller.NewAssetController(assetService, assetOnMarketService)
	assetEntryController = controller.NewAssetEntryController(entryService)
	assetOnMarketController = controller.NewAssetOnMarketController(assetOnMarketService, assetService)

	AvailableRoutes = []types.Route{
		{Method: "POST", Path: "/asset/add", Fn: assetController.AddUserAsset(), Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
		{Method: "GET", Path: "/asset/list", Fn: assetController.GetPaginatedUserAssets(), Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
		{Method: "GET", Path: "/asset/:id", Fn: assetController.GetUserAsset(), Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
		{Method: "GET", Path: "/asset/:code/realtime", Fn: assetOnMarketController.GetRealTimeAssetData(), Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
		{Method: "POST", Path: "/realtimeupdate", Fn: assetOnMarketController.UpdateAllAssetsOnMarket(), Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
		{Method: "POST", Path: "/asset/entry/add", Fn: assetEntryController.AddUserAssetEntry(), Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
		// {Method: "GET", Path: "/asset/entry/:id", Fn: GetUserAssetEntry, Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
	}
}
