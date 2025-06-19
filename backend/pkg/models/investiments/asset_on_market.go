package investiments

import (
	"easyinvesting/config"
	"easyinvesting/pkg/types"
	"encoding/json"
	"errors"
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
	AssetCode     string        `json:"asset_code" gorm:"not null"`
	AssetOnMarket AssetOnMarket `json:"asset_on_market" gorm:"foreignKey:AssetCode;references:Code;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Price         float64       `json:"price" gorm:"not null;default:0"`
	Date          string        `json:"date" gorm:"not null;uniqueIndex:idx_asset_date"`
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

	codesMerged := ""
	for i, asset := range assets {
		if i > 0 && i < len(assets)-1 {
			codesMerged += ","
		}
		codesMerged += asset.Code
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://brapi.dev/api/quote/"+codesMerged, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+config.BRAPI_TOKEN)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data Response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	if len(data.Results) == 0 {
		return errors.New("no real-time data found for the assets")
	}

	for _, quote := range data.Results {
		dailyAssetPrice := DailyAssetPrice{
			AssetCode:     quote.Symbol,
			AssetOnMarket: AssetOnMarket{Code: quote.Symbol},
			Price:         quote.RegularMarketPrice,
			Date:          time.Now().Format("2006-01-02"),
		}
		if err := config.DB.Create(&dailyAssetPrice).Error; err != nil {
			return err
		}
	}

	return nil
}
