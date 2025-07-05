package portfolio

import (
	"easyinvesting/config"
	"easyinvesting/pkg/api/v1/utils"
	"easyinvesting/pkg/model"
	"easyinvesting/pkg/types"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetRealTimeAssetData(c echo.Context) error {
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, types.JsonMap{"error": "Unauthorized"})
	}
	assetID := c.Param("id")
	if assetID == "" {
		return c.JSON(http.StatusBadRequest, types.JsonMap{"error": "Asset code is required"})
	}

	var asset model.Asset
	if err := config.DB().Select("code", "currency").Where("id = ? AND user_id = ?", assetID, claims.UserID).First(&asset).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, types.JsonMap{"error": "Asset not found"})
		}
		c.Logger().Errorf("Failed to fetch asset: %v", err)
		return c.JSON(http.StatusInternalServerError, types.JsonMap{"error": "Failed to fetch asset"})
	}

	var dailyPrice model.DailyAssetPrice
	if err := config.DB().Where("asset_code = ? AND date = ?", asset.Code, time.Now().Format("2006-01-02")).First(&dailyPrice).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, types.JsonMap{"error": "Daily price not found for today"})
		}
		c.Logger().Errorf("Failed to fetch daily price: %v", err)
		return c.JSON(http.StatusInternalServerError, types.JsonMap{"error": "Failed to fetch daily price"})
	}

	return c.JSON(http.StatusOK, types.JsonMap{
		"message":      "Real-time asset data retrieved successfully",
		"market_price": dailyPrice.Price,
	})
}

func UpdateRealTimeAssetsData(c echo.Context) error {
	_, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, types.JsonMap{"error": "Unauthorized"})
	}
	if err := model.UpdateAllAssetsOnMarket(); err != nil {
		c.Logger().Errorf("Failed to update real-time assets data: %v", err)
		return c.JSON(http.StatusInternalServerError, types.JsonMap{"error": "Failed to update real-time assets data"})
	}
	return c.JSON(http.StatusOK, types.JsonMap{
		"message": "Real-time assets data updated successfully",
	})
}
