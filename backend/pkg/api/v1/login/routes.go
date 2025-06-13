package login

import "easyinvesting/pkg/types"

var AvailableRoutes = []types.Route{
	{Method: "POST", Path: "/login", Fn: login},
	{Method: "POST", Path: "/signup", Fn: signup},
}
