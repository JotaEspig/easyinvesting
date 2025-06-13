package api

import (
	"easyinvesting/pkg/api/v1/login"
	"easyinvesting/pkg/api/v1/other"
	"easyinvesting/pkg/types"
)

var AllAvailableRoutes = []types.Route{}

func init() {
	AllAvailableRoutes = append(AllAvailableRoutes, other.AvailableRoutes...)
	AllAvailableRoutes = append(AllAvailableRoutes, login.AvailableRoutes...)

	for i := range AllAvailableRoutes {
		AllAvailableRoutes[i].Path = "/api/v1" + AllAvailableRoutes[i].Path
	}
}
