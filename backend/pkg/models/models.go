package models

import (
	"easyinvesting/config"
)

var allModels = []interface{}{
	&User{},
	&Asset{},
	&AssetEntry{},
	&AssetOnMarket{},
	&DailyAssetPrice{},
}

func Migrate() {
	db := config.DB()
	if err := db.AutoMigrate(
		allModels...,
	); err != nil {
		panic(err)
	}
}
