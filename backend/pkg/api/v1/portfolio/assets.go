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
