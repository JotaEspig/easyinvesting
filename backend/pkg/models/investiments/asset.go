package investiments

import (
	"easyinvesting/pkg/models/user"
	"easyinvesting/pkg/types"
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
	Name               string    `json:"name" gorm:"not null"`
	Code               string    `json:"code" gorm:"not null;unique"`
	AssetType          uint8     `json:"asset_type" gorm:"not null"`
	Currency           uint8     `json:"currency" gorm:"not null"`
	UserID             uint      `json:"user_id" gorm:"not null"`
	User               user.User `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CachedHoldAvgPrice float64   `json:"-" gorm:"default:0"`
	CachedHoldQuantity uint      `json:"-" gorm:"default:0"`
	CacheDate          time.Time `json:"-" gorm:"default:CURRENT_TIMESTAMP"`
}

func (a Asset) IsUserInputValid() bool {
	return a.Name != "" && a.Code != "" && a.AssetType == AssetTypeStock && a.Currency >= CurrencyBRL && a.Currency <= CurrencyUSD
}

func (a Asset) ToMap() types.JsonMap {
	return types.JsonMap{
		"id":             a.ID,
		"name":           a.Name,
		"code":           a.Code,
		"asset_type":     a.AssetType,
		"currency":       a.Currency,
		"user_id":        a.UserID,
		"hold_avg_price": a.CachedHoldAvgPrice,
		"hold_quantity":  a.CachedHoldQuantity,
	}
}
