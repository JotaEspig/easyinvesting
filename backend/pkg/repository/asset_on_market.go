package repository

import (
	"easyinvesting/pkg/model"

	"gorm.io/gorm"
)

type AssetOnMarketRepository interface {
	Save(assetOnMarket *model.AssetOnMarket) error
	GetAllAssetCodes() ([]string, error)
	EnsureAssetOnMarket(code string) (bool, error)
}

type assetOnMarketRepository struct {
	db *gorm.DB
}

func NewAssetOnMarketRepository(db *gorm.DB) AssetOnMarketRepository {
	return &assetOnMarketRepository{db: db}
}

func (r *assetOnMarketRepository) Save(assetOnMarket *model.AssetOnMarket) error {
	return r.db.Save(assetOnMarket).Error
}

func (r *assetOnMarketRepository) GetAllAssetCodes() ([]string, error) {
	var codes []string
	if err := r.db.Model(&model.AssetOnMarket{}).Pluck("code", &codes).Error; err != nil {
		return nil, err
	}
	return codes, nil
}

func (r *assetOnMarketRepository) EnsureAssetOnMarket(code string) (bool, error) {
	var assetOnMarket model.AssetOnMarket
	if err := r.db.Where("code = ?", code).First(&assetOnMarket).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return false, err
		}

		assetOnMarket = model.AssetOnMarket{Code: code}
		if err := r.Save(&assetOnMarket); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}
