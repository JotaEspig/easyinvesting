package repository

import (
	"easyinvesting/pkg/client"
	"easyinvesting/pkg/model"
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type DailyAssetPriceRepository interface {
	Save(dailyAssetPrice *model.DailyAssetPrice) error
	UpdateAllAssetsOnMarket(codes []string) error
	FindLatestPriceByCode(code string) (*model.DailyAssetPrice, error)
}

type dailyAssetPriceRepository struct {
	db *gorm.DB
}

func NewDailyAssetPriceRepository(db *gorm.DB) DailyAssetPriceRepository {
	return &dailyAssetPriceRepository{db: db}
}

func (r *dailyAssetPriceRepository) Save(dailyAssetPrice *model.DailyAssetPrice) error {
	return r.db.Save(dailyAssetPrice).Error
}

func (r *dailyAssetPriceRepository) UpdateAllAssetsOnMarket(codes []string) error {
	brApiClient := client.NewBrApi(&http.Client{})
	for _, code := range codes {
		var dayOfLastAssetUpdate string
		r.db.Model(&model.DailyAssetPrice{}).
			Select("MAX(date)").
			Where("asset_code = ?", code).
			Scan(&dayOfLastAssetUpdate)
		isToday := time.Now().Format("2006-01-02") == dayOfLastAssetUpdate
		if isToday {
			log.Printf("Asset %s already updated today, skipping...", code)
			continue
		}

		quote, err := brApiClient.GetQuote(code)
		if err != nil && err != client.BrApiErrNoResults {
			return fmt.Errorf("error fetching quote for %s: %w", code, err)
		}

		dailyAssetPrice := model.DailyAssetPrice{
			AssetCode: quote.Symbol,
			Price:     quote.RegularMarketPrice,
			Date:      time.Now().Format("2006-01-02"),
		}
		if err := r.Save(&dailyAssetPrice); err != nil {
			if err != gorm.ErrDuplicatedKey {
				return fmt.Errorf("error creating daily asset price for %s: %w", code, err)
			}
		}
	}

	return nil
}

func (r *dailyAssetPriceRepository) FindLatestPriceByCode(code string) (*model.DailyAssetPrice, error) {
	var dailyAssetPrice model.DailyAssetPrice
	err := r.db.Where("asset_code = ?", code).
		Order("date DESC").
		First(&dailyAssetPrice).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no daily asset price found for code %s", code)
		}
		return nil, fmt.Errorf("error fetching latest price for code %s: %w", code, err)
	}
	return &dailyAssetPrice, nil
}
