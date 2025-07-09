package service

import (
	"easyinvesting/config"
	"easyinvesting/pkg/dto"
	"easyinvesting/pkg/model"
	"easyinvesting/pkg/repository"
	"fmt"
)

type AssetEntryService interface {
	Save(assetDTO *dto.AssetEntryDTO, UserID uint) error
}

type assetEntryService struct {
	assetRepository      repository.AssetRepository
	assetEntryRepository repository.AssetEntryRepository
}

func NewAssetEntryService(
	assetRepository repository.AssetRepository,
	assetEntryRepository repository.AssetEntryRepository,
) AssetEntryService {
	return &assetEntryService{
		assetRepository:      assetRepository,
		assetEntryRepository: assetEntryRepository,
	}
}

func (s *assetEntryService) Save(assetEntryDTO *dto.AssetEntryDTO, UserID uint) error {
	if !assetEntryDTO.IsUserInputValid() {
		return fmt.Errorf("Invalid asset entry data: %v", assetEntryDTO.ToMap())
	}

	// Check if the asset exists for the user
	exists, err := s.assetRepository.DoesUserOwnAsset(assetEntryDTO.AssetID, UserID)
	if err != nil || !exists {
		return err
	}

	entry := &model.AssetEntry{
		Price:    assetEntryDTO.Price,
		Quantity: assetEntryDTO.Quantity,
		Type:     assetEntryDTO.Type,
		Date:     assetEntryDTO.Date,
		AssetID:  assetEntryDTO.AssetID,
	}

	tx := config.DB().Begin()
	if err := s.assetEntryRepository.SaveTx(tx, entry); err != nil {
		tx.Rollback()
		return err
	}

	// Update cached values in transaction
	if err := s.assetRepository.UpdateCachedValuesTx(tx, entry); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
