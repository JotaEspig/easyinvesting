package other

import (
	"easyinvesting/config"
	"easyinvesting/pkg/types"

	"github.com/labstack/echo/v4"
)

var AvailableRoutes []types.Route = []types.Route{
	{Path: "/helloauth", Method: types.MethodGET, Fn: helloAuth,
		Middlewares: []echo.MiddlewareFunc{config.JWTMiddleware()}},
}
