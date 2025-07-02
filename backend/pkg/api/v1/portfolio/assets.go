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

// curl command to test:
// curl -X POST -H "Content-Type: application/json" -d '{"name": "Petrobras", "code": "PETR3", "asset_type": 0}' http://localhost:8000/api/v1/asset/add
func AddUserAsset(c echo.Context) error {
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, types.JsonMap{"error": "Unauthorized"})
	}

	var a models.Asset
	if err := json.NewDecoder(c.Request().Body).Decode(&a); err != nil {
		c.Logger().Errorf("Failed to decode asset: %v", err.Error())
		return c.JSON(http.StatusBadRequest, types.JsonMap{
			"message": "some asset field may be missing or invalid",
		})
	}

	if !a.IsUserInputValid() {
		c.Logger().Errorf("Invalid asset input: %v", a)
		return c.JSON(http.StatusBadRequest, types.JsonMap{
			"message": "some asset field may be missing or invalid",
		})
	}

	err = ensureAssetOnMarket(c, &a)
	if err != nil {
		c.Logger().Errorf("Failed to ensure asset on market: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, types.JsonMap{
			"message": "Asset not on the market",
		})
	}

	a.UserID = claims.UserID
	if err := config.DB().Create(&a).Error; err != nil {
		c.Logger().Errorf("Failed to create asset: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, types.JsonMap{
			"message": "failed to create asset",
		})
	}
	c.Logger().Infof("Asset created successfully: %v", a)
	aDTO := dtos.AssetDTO{
		ID:                 a.ID,
		Code:               a.Code,
		AssetType:          a.AssetType,
		Currency:           a.Currency,
		UserID:             a.UserID,
		CachedHoldAvgPrice: a.CachedHoldAvgPrice,
		CachedHoldQuantity: a.CachedHoldQuantity,
	}
	return c.JSON(http.StatusCreated, types.JsonMap{
		"message": "asset created successfully",
		"asset":   aDTO.ToMap(),
	})
}

func GetUserAsset(c echo.Context) error {
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, types.JsonMap{"error": "Unauthorized"})
	}

	assetID := c.Param("id")
	var a models.Asset
	if err := config.DB().Where("id = ? AND user_id = ?", assetID, claims.UserID).First(&a).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, types.JsonMap{"message": "asset not found"})
		}
		c.Logger().Errorf("Failed to get asset: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, types.JsonMap{
			"message": "failed to get asset",
		})
	}

	c.Logger().Infof("Asset retrieved successfully: %v", a)
	aDTO := dtos.AssetDTO{
		ID:                 a.ID,
		Code:               a.Code,
		AssetType:          a.AssetType,
		Currency:           a.Currency,
		UserID:             a.UserID,
		CachedHoldAvgPrice: a.CachedHoldAvgPrice,
		CachedHoldQuantity: a.CachedHoldQuantity,
	}
	return c.JSON(http.StatusOK, types.JsonMap{
		"message": "asset retrieved successfully",
		"asset":   aDTO.ToMap(),
	})
}

func GetUserAssets(c echo.Context) error {
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, types.JsonMap{"error": "Unauthorized"})
	}

	var assets []models.Asset
	if err := config.DB().Where("user_id = ?", claims.UserID).Find(&assets).Error; err != nil {
		c.Logger().Errorf("Failed to get assets: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, types.JsonMap{
			"message": "failed to get assets",
		})
	}

	assetsMaps := make([]types.JsonMap, len(assets))
	for i, a := range assets {
		aDTO := dtos.AssetDTO{
			ID:                 a.ID,
			Code:               a.Code,
			AssetType:          a.AssetType,
			Currency:           a.Currency,
			UserID:             a.UserID,
			CachedHoldAvgPrice: a.CachedHoldAvgPrice,
			CachedHoldQuantity: a.CachedHoldQuantity,
		}
		assetsMaps[i] = aDTO.ToMap()
	}

	c.Logger().Infof("Assets retrieved successfully: %d assets", len(assets))
	return c.JSON(http.StatusOK, types.JsonMap{
		"message": "assets retrieved successfully",
		"assets":  assetsMaps,
	})
}

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
