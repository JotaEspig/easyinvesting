package investiments

import (
	"easyinvesting/config"
	"easyinvesting/pkg/types"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type AssetOnMarket struct {
	gorm.Model
	Code string `json:"code" gorm:"primary;not null;unique"`
}

func (AssetOnMarket) TableName() string {
	return "assets_on_market"
}

type DailyAssetPrice struct {
	gorm.Model
	AssetCode     string        `json:"asset_code" gorm:"not null;uniqueIndex:idx_code_date"`
	AssetOnMarket AssetOnMarket `json:"asset_on_market" gorm:"foreignKey:AssetCode;references:Code;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Price         float64       `json:"price" gorm:"not null;default:0"`
	Date          string        `json:"date" gorm:"not null;uniqueIndex:idx_code_date"`
}

func (dap DailyAssetPrice) ToMap() types.JsonMap {
	return map[string]interface{}{
		"asset_code": dap.AssetCode,
		"date":       dap.Date,
	}
}

type Quote struct {
	Symbol             string  `json:"symbol"`
	RegularMarketPrice float64 `json:"regularMarketPrice"`
}

type Response struct {
	Results []Quote `json:"results"`
}

func UpdateAllAssetsOnMarket() error {
	var assets []AssetOnMarket
	if err := config.DB.Select("code").Find(&assets).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	if len(assets) == 0 {
		return nil
	}

	client := &http.Client{}
	for _, asset := range assets {
		url := "https://brapi.dev/api/quote/" + asset.Code
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("error creating request for %s: %w", asset.Code, err)
		}
		req.Header.Set("Authorization", "Bearer "+config.BRAPI_TOKEN)

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error performing request for %s: %w", asset.Code, err)
		}
		defer resp.Body.Close()

		var data Response
		if resp.StatusCode != http.StatusOK {
			errorData := struct {
				Error   bool   `json:"error"`
				Message string `json:"message"`
			}{}

			if err := json.NewDecoder(resp.Body).Decode(&errorData); err != nil {
				return fmt.Errorf("error decoding error response for %s: %w", asset.Code, err)
			}
			return errors.New(errorData.Message)
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return fmt.Errorf("error decoding response for %s: %w", asset.Code, err)
		}

		if len(data.Results) == 0 {
			// This might happen if an asset code doesn't return data,
			// you can choose to skip or return an error.
			// For now, let's just log and continue for other assets.
			continue
		}

		quote := data.Results[0]
		dailyAssetPrice := DailyAssetPrice{
			AssetCode:     quote.Symbol,
			AssetOnMarket: AssetOnMarket{Code: quote.Symbol},
			Price:         quote.RegularMarketPrice,
			Date:          time.Now().Format("2006-01-02"),
		}
		if err := config.DB.Create(&dailyAssetPrice).Error; err != nil {
			if err != gorm.ErrDuplicatedKey {
				return fmt.Errorf("error creating daily asset price for %s: %w", asset.Code, err)
			}
		}
	}

	return nil
}
