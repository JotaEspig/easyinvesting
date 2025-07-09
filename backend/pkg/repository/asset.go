package repository

import (
	"easyinvesting/pkg/model"

	"gorm.io/gorm"
)

type AssetRepository interface {
	Save(asset *model.Asset) error
	FindByCodeAndUserID(code string, userID uint) (*model.Asset, error)
	GetPaginatedByUserID(userID uint, page, pageSize int) ([]*model.Asset, int64, error)
	DoesUserOwnAsset(id uint, userID uint) (bool, error)
	UpdateCachedValuesTx(tx *gorm.DB, entry *model.AssetEntry) error
}

type assetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepository{db: db}
}

func (r *assetRepository) Save(asset *model.Asset) error {
	return r.db.Save(asset).Error
}

func (r *assetRepository) FindByCodeAndUserID(code string, userID uint) (*model.Asset, error) {
	var asset model.Asset
	if err := r.db.Where("code = ? AND user_id = ?", code, userID).First(&asset).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *assetRepository) GetPaginatedByUserID(
	userID uint, page, pageSize int,
) ([]*model.Asset, int64, error) {
	var assets []*model.Asset
	var total int64

	query := r.db.Model(&model.Asset{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	if err := query.Find(&assets).Error; err != nil {
		return nil, 0, err
	}
	return assets, total, nil
}

func (r *assetRepository) DoesUserOwnAsset(id uint, userID uint) (bool, error) {
	var count int64
	if err := r.db.Model(&model.Asset{}).
		Where("id = ? AND user_id = ?", id, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *assetRepository) UpdateCachedValuesTx(tx *gorm.DB, entry *model.AssetEntry) error {
	if entry.Type == model.AssetEntryTypeBuy {
		return tx.Model(&model.Asset{}).Where("id = ?", entry.AssetID).
			UpdateColumns(map[string]interface{}{
				"cached_hold_avg_price": gorm.Expr(
					"((cached_hold_avg_price * cached_hold_quantity) + ?) / (cached_hold_quantity + ?)",
					entry.Price*float64(entry.Quantity), entry.Quantity,
				),
				"cached_hold_quantity": gorm.Expr("cached_hold_quantity + ?", entry.Quantity),
				"cache_date":           gorm.Expr("CURRENT_TIMESTAMP"),
			}).Error
	} else if entry.Type == model.AssetEntryTypeSell {
		return tx.Model(&model.Asset{}).Where("id = ?", entry.AssetID).
			UpdateColumns(map[string]interface{}{
				"cached_hold_avg_price": gorm.Expr(
					"((cached_hold_avg_price * cached_hold_quantity) - ?) / (cached_hold_quantity - ?)",
					entry.Price*float64(entry.Quantity), entry.Quantity,
				),
				"cached_hold_quantity": gorm.Expr("cached_hold_quantity - ?", entry.Quantity),
				"cache_date":           gorm.Expr("CURRENT_TIMESTAMP"),
			}).Error
	}
	return nil
}
