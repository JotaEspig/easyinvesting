package models

import (
	"easyinvesting/pkg/models/investiments"
	"easyinvesting/pkg/models/user"
)

var AllModels = []interface{}{
	&user.User{},
	&investiments.Asset{},
	&investiments.AssetEntry{},
}
