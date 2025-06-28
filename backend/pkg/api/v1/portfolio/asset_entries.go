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
// curl -X POST -H "Content-Type: application/json" -d '{"price": 10.0, "quantity": 2, "type": 0, "date": "2025-06-19T10:00:00Z", "asset_id": 3}' http://localhost:8000/api/v1/asset/entry/add
func AddUserAssetEntry(c echo.Context) error {
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, types.JsonMap{"error": "Unauthorized"})
	}

	var entry investiments.AssetEntry
	if err := json.NewDecoder(c.Request().Body).Decode(&entry); err != nil {
		c.Logger().Errorf("Failed to decode asset entry: %v", err.Error())
		return c.JSON(http.StatusBadRequest, types.JsonMap{
			"message": "some asset entry field may be missing or invalid",
		})
	}

	if !entry.IsUserInputValid() {
		c.Logger().Errorf("Invalid asset entry input: %v", entry)
		return c.JSON(http.StatusBadRequest, types.JsonMap{
			"message": "some asset entry field may be missing or invalid",
		})
	}

	// Verify if user has this asset
	var exists bool
	err = config.DB.Model(&investiments.Asset{}).
		Select("1").
		Where("id = ? AND user_id = ?", entry.AssetID, claims.UserID).
		Limit(1).
		Find(&exists).Error
	if err != nil || !exists {
		c.Logger().Errorf("Asset not found for user: %v", err.Error())
		return c.JSON(http.StatusNotFound, types.JsonMap{
			"message": "asset not found for user",
		})
	}

	// Verify if it's a sell entry and if the user has enough quantity
	if entry.Type == investiments.AssetEntryTypeSell {
		var asset investiments.Asset
		err = config.DB.Model(&investiments.Asset{}).
			Select("cached_hold_quantity").
			Where("id = ? AND user_id = ?", entry.AssetID, claims.UserID).
			First(&asset).Error
		if err != nil {
			c.Logger().Errorf("Failed to fetch asset for sell entry: %v", err.Error())
			return c.JSON(http.StatusInternalServerError, types.JsonMap{
				"message": "failed to fetch asset for sell entry",
			})
		}
		if asset.CachedHoldQuantity < entry.Quantity {
			c.Logger().Errorf("Insufficient quantity for sell entry: %v", asset.CachedHoldQuantity)
			return c.JSON(http.StatusBadRequest, types.JsonMap{
				"message": "insufficient quantity for sell entry",
			})
		}
	}

	tx := config.DB.Begin()
	if err := tx.Create(&entry).Error; err != nil {
		defer tx.Rollback()
		c.Logger().Errorf("Failed to create asset entry: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, types.JsonMap{
			"message": "failed to create asset entry",
		})
	}

	// Update asset's cached hold average price and quantity
	if entry.Type == investiments.AssetEntryTypeBuy {
		err = tx.Model(&investiments.Asset{}).Where("id = ?", entry.AssetID).
			UpdateColumns(map[string]interface{}{
				"cached_hold_avg_price": gorm.Expr(
					"((cached_hold_avg_price * cached_hold_quantity) + ?) / (cached_hold_quantity + ?)",
					entry.Price*float64(entry.Quantity), entry.Quantity,
				),
				"cached_hold_quantity": gorm.Expr("cached_hold_quantity + ?", entry.Quantity),
				"cache_date":           gorm.Expr("CURRENT_TIMESTAMP"),
			}).Error
	} else if entry.Type == investiments.AssetEntryTypeSell {
		err = tx.Model(&investiments.Asset{}).Where("id = ?", entry.AssetID).
			UpdateColumns(map[string]interface{}{
				"cached_hold_avg_price": gorm.Expr(
					"((cached_hold_avg_price * cached_hold_quantity) - ?) / (cached_hold_quantity - ?)",
					entry.Price*float64(entry.Quantity), entry.Quantity,
				),
				"cached_hold_quantity": gorm.Expr("cached_hold_quantity - ?", entry.Quantity),
				"cache_date":           gorm.Expr("CURRENT_TIMESTAMP"),
			}).Error
	}
	if err != nil {
		defer tx.Rollback()
		c.Logger().Errorf("Failed to update asset cached values: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, types.JsonMap{
			"message": "failed to update asset cached values",
		})
	}

	tx.Commit()
	c.Logger().Infof("Asset entry created successfully: %v", entry)
	return c.JSON(http.StatusCreated, types.JsonMap{
		"message": "asset entry created successfully",
		"entry":   entry.ToMap(),
	})
}

func GetUserAssetEntry(c echo.Context) error {
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, types.JsonMap{"error": "Unauthorized"})
	}

	entryID := c.Param("id")
	var entry investiments.AssetEntry
	err = config.DB.Table("asset_entries").
		Select("asset_entries.*").
		Joins("INNER JOIN assets ON asset_entries.asset_id = assets.id").
		Where("asset_entries.id = ? AND assets.user_id = ?", entryID, claims.UserID).
		First(&entry).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, types.JsonMap{"message": "asset entry not found"})
		}
		c.Logger().Errorf("Database error while fetching asset entry: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, types.JsonMap{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, types.JsonMap{
		"message": "asset entry retrieved successfully",
		"entry":   entry.ToMap(),
	})
}
