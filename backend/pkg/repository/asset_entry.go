package repository

import (
	"easyinvesting/pkg/model"

	"gorm.io/gorm"
)

type AssetEntryRepository interface {
	Save(entry *model.AssetEntry) error
	SaveTx(tx *gorm.DB, entry *model.AssetEntry) error
	FindByIDAndUserID(id, userID uint) (*model.AssetEntry, error)
}

type assetEntryRepository struct {
	db *gorm.DB
}

func NewAssetEntryRepository(db *gorm.DB) AssetEntryRepository {
	return &assetEntryRepository{db: db}
}

func (r *assetEntryRepository) Save(entry *model.AssetEntry) error {
	return r.db.Save(entry).Error
}

func (r *assetEntryRepository) SaveTx(tx *gorm.DB, entry *model.AssetEntry) error {
	return tx.Save(entry).Error
}

func (r *assetEntryRepository) FindByIDAndUserID(id, userID uint) (*model.AssetEntry, error) {
	var entry model.AssetEntry
	if err := r.db.Model(&entry).Where("id = ? AND user_id = ?", id, userID).First(&entry).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}
