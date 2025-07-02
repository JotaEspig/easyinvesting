package dtos

const (
	AssetTypeStock uint8 = 0
)

const (
	CurrencyBRL uint8 = 0
	CurrencyUSD uint8 = 1
)

type AssetDTO struct {
	ID                 uint    `json:"id"`
	Code               string  `json:"code"`
	AssetType          uint8   `json:"asset_type"`
	Currency           uint8   `json:"currency"`
	UserID             uint    `json:"user_id"`
	CachedHoldAvgPrice float64 `json:"hold_avg_price"`
	CachedHoldQuantity uint    `json:"hold_quantity"`
}

func (a AssetDTO) IsUserInputValid() bool {
	return a.Code != "" && a.AssetType == AssetTypeStock && a.Currency >= CurrencyBRL && a.Currency <= CurrencyUSD
}

func (a AssetDTO) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":             a.ID,
		"code":           a.Code,
		"asset_type":     a.AssetType,
		"currency":       a.Currency,
		"user_id":        a.UserID,
		"hold_avg_price": a.CachedHoldAvgPrice,
		"hold_quantity":  a.CachedHoldQuantity,
	}
}
