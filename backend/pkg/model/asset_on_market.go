package model

import (
	"gorm.io/gorm"
)

type AssetOnMarket struct {
	gorm.Model
	Code string `json:"code" gorm:"primary;not null;unique"`
}

func (AssetOnMarket) TableName() string {
	return "assets_on_market"
}
