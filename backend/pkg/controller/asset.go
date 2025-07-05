package controller

import (
	"easyinvesting/config"
	"easyinvesting/pkg/client"
	"easyinvesting/pkg/controller/utils"
	"easyinvesting/pkg/dto"
	"easyinvesting/pkg/model"
	"easyinvesting/pkg/service"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AssetController struct {
	assetService service.AssetService
}

func NewAssetController(assetService service.AssetService) *AssetController {
	return &AssetController{assetService: assetService}
}

func (controller AssetController) AddUserAsset() echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := utils.GetClaimsFromContext(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		}

		var asset dto.AssetDTO
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

		// TODO use this instead: err = controller.assetService.EnsureAssetOnMarket(c, &asset)
		if err := ensureAssetOnMarket(c, &asset); err != nil {
			c.Logger().Errorf("Failed to ensure asset on market: %v", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "failed to ensure asset on market",
			})
		}

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
		// if pageStr == "" || pageSizeStr == "" {
		// 	return c.JSON(http.StatusBadRequest, map[string]string{"message": "Page and size parameters are required"})
		// }
		// This below is temporary, in future, use the above one
		if pageStr == "" {
			pageStr = "1"
		}
		if pageSizeStr == "" {
			pageSizeStr = "10"
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

// CHANGE THIS IN FUTURE
func ensureAssetOnMarket(c echo.Context, a *dto.AssetDTO) error {
	var assetOnMarket model.AssetOnMarket
	if err := config.DB().Where("code = ?", a.Code).First(&assetOnMarket).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.Logger().Errorf("Failed to check asset on market: %v", err.Error())
			return err
		}

		brApiClient := client.NewBrApi(&http.Client{})
		quote, err := brApiClient.GetQuote(a.Code)

		if err != nil {
			if err == client.BrApiErrNoResults {
				c.Logger().Errorf("Asset not found on market: %s", a.Code)
				return echo.NewHTTPError(http.StatusNotFound, "Asset not found on market")
			}
			c.Logger().Errorf("Failed to fetch asset quote: %v", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch asset quote")
		}

		assetOnMarket = model.AssetOnMarket{Code: a.Code}
		if err := config.DB().Create(&assetOnMarket).Error; err != nil {
			c.Logger().Errorf("Failed to create asset on market: %v", err.Error())
			return err
		}
		c.Logger().Infof("Asset on market created: %s", a.Code)

		// create daily asset price
		dailyAssetPrice := model.DailyAssetPrice{
			AssetCode:     quote.Symbol,
			AssetOnMarket: model.AssetOnMarket{Code: quote.Symbol},
			Price:         quote.RegularMarketPrice,
			Date:          time.Now().Format("2006-01-02"),
		}
		if err := config.DB().Create(&dailyAssetPrice).Error; err != nil {
			if err != gorm.ErrDuplicatedKey {
				return fmt.Errorf("error creating daily asset price for %s: %w", a.Code, err)
			}
		}
	}

	// a.AssetOnMarket = assetOnMarket
	return nil
}
