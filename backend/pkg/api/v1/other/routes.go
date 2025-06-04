package other

import (
	"easyinvesting/pkg/types"
)

var AvailableRoutes []types.Route = []types.Route{
	{Path: "/api/hello", Method: types.MethodGET, Fn: hello},
}
