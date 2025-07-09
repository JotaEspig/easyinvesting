package dto

type DailyAssetPriceDTO struct {
	AssetCode string  `json:"asset_code"`
	Price     float64 `json:"price"`
	Date      string  `json:"date"`
}

func (dap DailyAssetPriceDTO) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"asset_code": dap.AssetCode,
		"date":       dap.Date,
	}
}
