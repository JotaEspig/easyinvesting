package model

import (
	"easyinvesting/config"
	"easyinvesting/pkg/client"
	"easyinvesting/pkg/types"
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

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
	if err := config.DB().Select("code").Find(&assets).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	if len(assets) == 0 {
		return nil
	}

	brApiClient := client.NewBrApi(&http.Client{})
	for _, asset := range assets {
		var dayOfLastAssetUpdate string
		config.DB().Model(&DailyAssetPrice{}).
			Select("MAX(date)").
			Where("asset_code = ?", asset.Code).
			Scan(&dayOfLastAssetUpdate)
		isToday := time.Now().Format("2006-01-02") == dayOfLastAssetUpdate
		if isToday {
			log.Printf("Asset %s already updated today, skipping...", asset.Code)
			continue
		}

		quote, err := brApiClient.GetQuote(asset.Code)
		if err != nil && err != client.BrApiErrNoResults {
			return fmt.Errorf("error fetching quote for %s: %w", asset.Code, err)
		}

		dailyAssetPrice := DailyAssetPrice{
			AssetCode:     quote.Symbol,
			AssetOnMarket: AssetOnMarket{Code: quote.Symbol},
			Price:         quote.RegularMarketPrice,
			Date:          time.Now().Format("2006-01-02"),
		}
		if err := config.DB().Create(&dailyAssetPrice).Error; err != nil {
			if err != gorm.ErrDuplicatedKey {
				return fmt.Errorf("error creating daily asset price for %s: %w", asset.Code, err)
			}
		}
	}

	return nil
}
