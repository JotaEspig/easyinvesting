package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	AssetTypeStock uint8 = 0
)

const (
	CurrencyBRL uint8 = 0
	CurrencyUSD uint8 = 1
)

type Asset struct {
	gorm.Model
	Code               string        `gorm:"not null"`
	AssetOnMarket      AssetOnMarket `gorm:"foreignKey:Code;references:Code;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AssetType          uint8         `gorm:"not null"`
	Currency           uint8         `gorm:"not null"`
	UserID             uint          `gorm:"not null"`
	User               User          `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CachedHoldAvgPrice float64       `gorm:"default:0"`
	CachedHoldQuantity uint          `gorm:"default:0"`
	CacheDate          time.Time     `gorm:"default:CURRENT_TIMESTAMP"`
}

func (a Asset) IsUserInputValid() bool {
	return a.Code != "" && a.AssetType == AssetTypeStock && a.Currency >= CurrencyBRL && a.Currency <= CurrencyUSD
}
