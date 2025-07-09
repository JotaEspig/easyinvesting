package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	AssetEntryTypeBuy  uint8 = 0
	AssetEntryTypeSell uint8 = 1
)

type AssetEntry struct {
	gorm.Model
	Price    float64   `gorm:"not null"`
	Quantity uint      `gorm:"not null"`
	Type     uint8     `gorm:"not null"`
	Date     time.Time `gorm:"not null"` // RFC 3339
	AssetID  uint      `gorm:"not null"`
	Asset    Asset     `gorm:"foreignKey:AssetID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
