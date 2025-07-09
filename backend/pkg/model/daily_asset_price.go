package model

import (
	"gorm.io/gorm"
)

type DailyAssetPrice struct {
	gorm.Model
	AssetCode     string        `gorm:"not null;uniqueIndex:idx_code_date"`
	AssetOnMarket AssetOnMarket `gorm:"foreignKey:AssetCode;references:Code;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Price         float64       `gorm:"not null;default:0"`
	Date          string        `gorm:"not null;uniqueIndex:idx_code_date"`
}
