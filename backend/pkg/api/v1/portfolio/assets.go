package portfolio

import (
	"easyinvesting/config"
	"easyinvesting/pkg/api/v1/utils"
	"easyinvesting/pkg/models/investiments"
	"easyinvesting/pkg/types"
	"encoding/json"
	"net/http"

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

	var a investiments.Asset
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
	if err := config.DB.Create(&a).Error; err != nil {
		c.Logger().Errorf("Failed to create asset: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, types.JsonMap{
			"message": "failed to create asset",
		})
	}
	c.Logger().Infof("Asset created successfully: %v", a)
	return c.JSON(http.StatusCreated, types.JsonMap{
		"message": "asset created successfully",
		"asset":   a.ToMap(),
	})
}

func GetUserAsset(c echo.Context) error {
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, types.JsonMap{"error": "Unauthorized"})
	}

	assetID := c.Param("id")
	var a investiments.Asset
	if err := config.DB.Where("id = ? AND user_id = ?", assetID, claims.UserID).First(&a).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, types.JsonMap{"message": "asset not found"})
		}
		c.Logger().Errorf("Failed to get asset: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, types.JsonMap{
			"message": "failed to get asset",
		})
	}

	c.Logger().Infof("Asset retrieved successfully: %v", a)
	return c.JSON(http.StatusOK, types.JsonMap{
		"message": "asset retrieved successfully",
		"asset":   a.ToMap(),
	})
}

func GetUserAssets(c echo.Context) error {
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, types.JsonMap{"error": "Unauthorized"})
	}

	var assets []investiments.Asset
	if err := config.DB.Where("user_id = ?", claims.UserID).Find(&assets).Error; err != nil {
		c.Logger().Errorf("Failed to get assets: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, types.JsonMap{
			"message": "failed to get assets",
		})
	}

	assetsMaps := make([]types.JsonMap, len(assets))
	for i, a := range assets {
		assetsMaps[i] = a.ToMap()
	}

	c.Logger().Infof("Assets retrieved successfully: %d assets", len(assets))
	return c.JSON(http.StatusOK, types.JsonMap{
		"message": "assets retrieved successfully",
		"assets":  assetsMaps,
	})
}

func ensureAssetOnMarket(c echo.Context, a *investiments.Asset) error {
	var assetOnMarket investiments.AssetOnMarket
	if err := config.DB.Where("code = ?", a.Code).First(&assetOnMarket).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// check if the asset exists in the market
			client := &http.Client{}
			req, err := http.NewRequest("GET", "https://brapi.dev/api/quote/"+a.Code, nil)
			if err != nil {
				c.Logger().Errorf("Failed to create request: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create request")
			}
			req.Header.Set("Authorization", "Bearer "+config.BRAPI_TOKEN)
			resp, err := client.Do(req)
			if err != nil {
				c.Logger().Errorf("Failed to fetch real-time data: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch real-time data")
			}
			defer resp.Body.Close()

			var data investiments.Response
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				c.Logger().Errorf("Failed to decode response: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode response")
			}

			if len(data.Results) == 0 {
				c.Logger().Errorf("No real-time data found for the asset: %s", a.Code)
				return c.JSON(http.StatusNotFound, types.JsonMap{"error": "No real-time data found for the asset"})
			}

			assetOnMarket = investiments.AssetOnMarket{Code: a.Code}
			if err := config.DB.Create(&assetOnMarket).Error; err != nil {
				c.Logger().Errorf("Failed to create asset on market: %v", err.Error())
				return err
			}
			c.Logger().Infof("Asset on market created: %s", a.Code)
		} else {
			c.Logger().Errorf("Failed to check asset on market: %v", err.Error())
			return err
		}
	}

	a.AssetOnMarket = assetOnMarket
	return nil
}
