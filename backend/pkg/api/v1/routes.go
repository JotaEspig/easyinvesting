package api

import (
	"easyinvesting/pkg/api/v1/login"
	"easyinvesting/pkg/api/v1/other"
	"easyinvesting/pkg/api/v1/portfolio"
	"easyinvesting/pkg/types"
)

var AllAvailableRoutes = []types.Route{}

func init() {
	AllAvailableRoutes = append(AllAvailableRoutes, other.AvailableRoutes...)
	AllAvailableRoutes = append(AllAvailableRoutes, login.AvailableRoutes...)
	AllAvailableRoutes = append(AllAvailableRoutes, portfolio.AvailableRoutes...)
}
