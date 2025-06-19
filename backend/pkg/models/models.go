package models

import (
	"easyinvesting/config"
	"easyinvesting/pkg/models/investiments"
	"easyinvesting/pkg/models/user"
)

var allModels = []interface{}{
	&user.User{},
	&investiments.Asset{},
	&investiments.AssetEntry{},
	&investiments.AssetOnMarket{},
	&investiments.DailyAssetPrice{},
}

func Migrate() {
	config.InitDB()
	if err := config.DB.AutoMigrate(
		allModels...,
	); err != nil {
		panic(err)
	}
}
