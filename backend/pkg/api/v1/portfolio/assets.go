package portfolio

import (
	"easyinvesting/config"
	"easyinvesting/pkg/api/v1/utils"
	"easyinvesting/pkg/clients"
	"easyinvesting/pkg/dtos"
	"easyinvesting/pkg/models"
	"easyinvesting/pkg/types"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ensureAssetOnMarket(c echo.Context, a *models.Asset) error {
	var assetOnMarket models.AssetOnMarket
	if err := config.DB().Where("code = ?", a.Code).First(&assetOnMarket).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.Logger().Errorf("Failed to check asset on market: %v", err.Error())
			return err
		}

		client := clients.NewBrApi(&http.Client{})
		quote, err := client.GetQuote(a.Code)

		if err != nil {
			if err == clients.BrApiErrNoResults {
				c.Logger().Errorf("Asset not found on market: %s", a.Code)
				return echo.NewHTTPError(http.StatusNotFound, "Asset not found on market")
			}
			c.Logger().Errorf("Failed to fetch asset quote: %v", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch asset quote")
		}

		assetOnMarket = models.AssetOnMarket{Code: a.Code}
		if err := config.DB().Create(&assetOnMarket).Error; err != nil {
			c.Logger().Errorf("Failed to create asset on market: %v", err.Error())
			return err
		}
		c.Logger().Infof("Asset on market created: %s", a.Code)

		// create daily asset price
		dailyAssetPrice := models.DailyAssetPrice{
			AssetCode:     quote.Symbol,
			AssetOnMarket: models.AssetOnMarket{Code: quote.Symbol},
			Price:         quote.RegularMarketPrice,
			Date:          time.Now().Format("2006-01-02"),
		}
		if err := config.DB().Create(&dailyAssetPrice).Error; err != nil {
			if err != gorm.ErrDuplicatedKey {
				return fmt.Errorf("error creating daily asset price for %s: %w", a.Code, err)
			}
		}
	}

	a.AssetOnMarket = assetOnMarket
	return nil
}
