package controller

import (
	"easyinvesting/pkg/controller/utils"
	"easyinvesting/pkg/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AssetOnMarketController struct {
	assetService         service.AssetService
	assetOnMarketService service.AssetOnMarketService
}

func NewAssetOnMarketController(
	assetOnMarketService service.AssetOnMarketService,
	assetService service.AssetService,
) *AssetOnMarketController {
	return &AssetOnMarketController{
		assetOnMarketService: assetOnMarketService,
		assetService:         assetService,
	}
}

func (c *AssetOnMarketController) GetRealTimeAssetData() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		claims, err := utils.GetClaimsFromContext(ctx)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Unauthorized"})
		}

		code := ctx.Param("code")
		if code == "" {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Asset code is required"})
		}

		exists, err := c.assetService.DoesUserOwnAsset(code, claims.UserID)
		if err != nil || !exists {
			if err != nil {
				ctx.Logger().Errorf("Failed to check asset ownership: %v", err)
				return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check asset ownership"})
			}
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Asset not found or not owned by user"})
		}

		data, err := c.assetOnMarketService.GetRealTimeAssetData(code)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get real-time asset data"})
		}

		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"market_price": data.Price,
		})
	}
}

func (c *AssetOnMarketController) UpdateAllAssetsOnMarket() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		_, err := utils.GetClaimsFromContext(ctx)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		if err := c.assetOnMarketService.UpdateAllAssetsOnMarket(); err != nil {
			ctx.Logger().Errorf("Failed to update all assets on market: %v", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update all assets on market"})
		}
		return ctx.JSON(http.StatusOK, map[string]string{"message": "All assets on market updated successfully"})
	}
}
