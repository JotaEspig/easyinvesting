package controllers

import (
	"easyinvesting/pkg/controllers/utils"
	"easyinvesting/pkg/dtos"
	"easyinvesting/pkg/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AssetController struct {
	assetService services.AssetService
}

func NewAssetController(assetService services.AssetService) *AssetController {
	return &AssetController{assetService: assetService}
}

func (controller AssetController) AddUserAsset() echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := utils.GetClaimsFromContext(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		}

		var asset dtos.AssetDTO
		if err := c.Bind(&asset); err != nil {
			c.Logger().Errorf("Failed to decode asset: %v", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "some asset field may be missing or invalid",
			})
		}

		if !asset.IsUserInputValid() {
			c.Logger().Errorf("Invalid asset input: %v", asset)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "some asset field may be missing or invalid",
			})
		}

		// TODO: err = controller.assetService.EnsureAssetOnMarket(c, &asset)

		asset.UserID = claims.UserID
		if err := controller.assetService.Save(&asset); err != nil {
			c.Logger().Errorf("Failed to create asset: %v", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "failed to create asset",
			})
		}

		c.Logger().Infof("Asset created successfully: %v", asset)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Asset created successfully",
			"asset":   asset.ToMap(),
		})
	}
}

func (controller AssetController) GetUserAsset() echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := utils.GetClaimsFromContext(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		}

		assetID := c.Param("id")
		if assetID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Asset ID is required"})
		}

		asset, err := controller.assetService.FindByCodeAndUserID(assetID, claims.UserID)
		if err != nil {
			c.Logger().Errorf("Failed to find asset: %v", err.Error())
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Asset not found"})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"asset": asset.ToMap(),
		})
	}
}

func (controller AssetController) GetPaginatedUserAssets() echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := utils.GetClaimsFromContext(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		}

		pageStr := c.QueryParam("page")
		pageSizeStr := c.QueryParam("size")
		if pageStr == "" || pageSizeStr == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Page and size parameters are required"})
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			c.Logger().Errorf("Invalid page parameter: %v", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid page parameter"})
		}
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			c.Logger().Errorf("Invalid size parameter: %v", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid size parameter"})
		}

		if page < 1 || pageSize < 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid pagination parameters"})
		}

		assets, total, err := controller.assetService.GetPaginatedByUserID(claims.UserID, page, pageSize)
		if err != nil {
			c.Logger().Errorf("Failed to get paginated assets: %v", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get assets"})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"assets": assets,
			"total":  total,
			"page":   page,
			"size":   pageSize,
		})
	}
}
