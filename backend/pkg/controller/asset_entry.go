package controller

import (
	"easyinvesting/pkg/controller/utils"
	"easyinvesting/pkg/dto"
	"easyinvesting/pkg/service"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AssetEntryController struct {
	assetEntryService service.AssetEntryService
}

func NewAssetEntryController(assetEntryService service.AssetEntryService) *AssetEntryController {
	return &AssetEntryController{assetEntryService: assetEntryService}
}

func (controller *AssetEntryController) AddUserAssetEntry() echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := utils.GetClaimsFromContext(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		}

		var assetEntry dto.AssetEntryDTO
		if err := c.Bind(&assetEntry); err != nil {
			c.Logger().Errorf("Failed to decode asset entry: %v", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "some asset entry field may be missing or invalid",
			})
		}

		if !assetEntry.IsUserInputValid() {
			c.Logger().Errorf("Invalid asset entry input: %v", assetEntry)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "some asset entry field may be missing or invalid",
			})
		}

		if err := controller.assetEntryService.Save(&assetEntry, claims.UserID); err != nil {
			c.Logger().Errorf("Failed to create asset entry: %v", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": fmt.Sprintf("failed to create asset entry: %v", err.Error()),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Asset entry created successfully"})
	}
}
