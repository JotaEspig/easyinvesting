package portfolio

import (
	"easyinvesting/config"
	"easyinvesting/pkg/types"

	"github.com/labstack/echo/v4"
)

var AvailableRoutes = []types.Route{
	{Method: "POST", Path: "/asset/add", Fn: AddUserAsset, Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
	{Method: "GET", Path: "/asset/list", Fn: GetUserAssets, Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
	{Method: "GET", Path: "/asset/:id", Fn: GetUserAsset, Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
	{Method: "POST", Path: "/asset/entry/add", Fn: AddUserAssetEntry, Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
	{Method: "GET", Path: "/asset/entry/:id", Fn: GetUserAssetEntry, Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
}
