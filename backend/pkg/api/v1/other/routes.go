package other

import (
	"easyinvesting/config"
	"easyinvesting/pkg/types"

	"github.com/labstack/echo/v4"
)

var AvailableRoutes []types.Route = []types.Route{
	{Path: "/hello", Method: types.MethodGET, Fn: hello},
	{Path: "/helloauth", Method: types.MethodGET, Fn: helloAuth,
		Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
}
