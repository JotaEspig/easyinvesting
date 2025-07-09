package dto

import "time"

const (
	AssetEntryTypeBuy  uint8 = 0
	AssetEntryTypeSell uint8 = 1
)

type AssetEntryDTO struct {
	ID       uint      `json:"id"`
	Price    float64   `json:"price"`
	Quantity uint      `json:"quantity"`
	Type     uint8     `json:"type"`
	Date     time.Time `json:"date"` // RFC 3339
	AssetID  uint      `json:"asset_id"`
}

func (ae AssetEntryDTO) IsUserInputValid() bool {
	return ae.Price > 0 && ae.Quantity > 0 && (ae.Type == AssetEntryTypeBuy || ae.Type == AssetEntryTypeSell) && !ae.Date.IsZero() && ae.AssetID > 0
}

func (ae AssetEntryDTO) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":       ae.ID,
		"price":    ae.Price,
		"quantity": ae.Quantity,
		"type":     ae.Type,
		"date":     ae.Date.Format(time.RFC3339),
		"asset_id": ae.AssetID,
	}
}
