package repositories

import (
	"easyinvesting/pkg/models"

	"gorm.io/gorm"
)

type AssetRepository interface {
	Save(asset *models.Asset) error
	FindByCodeAndUserID(code string, userID uint) (*models.Asset, error)
	GetPaginatedByUserID(userID uint, page, pageSize int) ([]*models.Asset, int64, error)
}

type assetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepository{db: db}
}

func (r *assetRepository) Save(asset *models.Asset) error {
	return r.db.Save(asset).Error
}

func (r *assetRepository) FindByCodeAndUserID(code string, userID uint) (*models.Asset, error) {
	var asset models.Asset
	if err := r.db.Where("code = ? AND user_id = ?", code, userID).First(&asset).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *assetRepository) GetPaginatedByUserID(
	userID uint, page, pageSize int,
) ([]*models.Asset, int64, error) {
	var assets []*models.Asset
	var total int64

	query := r.db.Model(&models.Asset{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	if err := query.Find(&assets).Error; err != nil {
		return nil, 0, err
	}
	return assets, total, nil
}
