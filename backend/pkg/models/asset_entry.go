package models

import (
	"easyinvesting/pkg/types"
	"time"

	"gorm.io/gorm"
)

const (
	AssetEntryTypeBuy  uint8 = 0
	AssetEntryTypeSell uint8 = 1
)

type AssetEntry struct {
	gorm.Model
	Price    float64   `json:"price" gorm:"not null"`
	Quantity uint      `json:"quantity" gorm:"not null"`
	Type     uint8     `json:"type" gorm:"not null"`
	Date     time.Time `json:"date" gorm:"not null"` // RFC 3339
	AssetID  uint      `json:"asset_id" gorm:"not null"`
	Asset    Asset     `json:"asset" gorm:"foreignKey:AssetID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ae AssetEntry) IsUserInputValid() bool {
	return ae.Price > 0 && ae.Quantity > 0 && (ae.Type == AssetEntryTypeBuy || ae.Type == AssetEntryTypeSell) && !ae.Date.IsZero() && ae.AssetID > 0
}

func (ae AssetEntry) ToMap() types.JsonMap {
	return types.JsonMap{
		"id":       ae.ID,
		"price":    ae.Price,
		"quantity": ae.Quantity,
		"type":     ae.Type,
		"date":     ae.Date.Format(time.RFC3339),
		"asset_id": ae.AssetID,
	}
}
