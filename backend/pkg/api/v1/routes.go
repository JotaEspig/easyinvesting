package api

import (
	"easyinvesting/pkg/api/v1/other"
	"easyinvesting/pkg/types"
)

var AllAvailableRoutes = []types.Route{}

func init() {
	AllAvailableRoutes = append(AllAvailableRoutes, other.AvailableRoutes...)
}
