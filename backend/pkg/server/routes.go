package server

import (
	api "easyinvesting/pkg/api/v1"
	"easyinvesting/pkg/types"

	"github.com/labstack/echo/v4"
)

func (s *Server) setRoute(group *echo.Group, route types.Route) {
	gp := group.Group(route.Path, route.Middlewares...)
	switch route.Method {
	case types.MethodGET:
		gp.GET("", route.Fn)
	case types.MethodPOST:
		gp.POST("", route.Fn)
	case types.MethodPUT:
		gp.PUT("", route.Fn)
	case types.MethodDELETE:
		gp.DELETE("", route.Fn)
	case types.MethodOPTIONS:
		gp.OPTIONS("", route.Fn)
	}
}

func (s *Server) setRoutes() {
	group := s.echo.Group("/api/v1")
	for _, route := range api.AllAvailableRoutes {
		s.setRoute(group, route)
	}
}
