package repository

import (
	"easyinvesting/pkg/model"

	"gorm.io/gorm"
)

type AssetRepository interface {
	Save(asset *model.Asset) error
	FindByCodeAndUserID(code string, userID uint) (*model.Asset, error)
	GetPaginatedByUserID(userID uint, page, pageSize int) ([]*model.Asset, int64, error)
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
