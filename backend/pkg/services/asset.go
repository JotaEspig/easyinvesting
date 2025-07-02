package services

import (
	"easyinvesting/pkg/dtos"
	"easyinvesting/pkg/models"
	"easyinvesting/pkg/repositories"
	"fmt"
)

type AssetService interface {
	Save(assetDTO *dtos.AssetDTO) error
	FindByCodeAndUserID(code string, userID uint) (*dtos.AssetDTO, error)
	GetPaginatedByUserID(userID uint, page, pageSize int) ([]*dtos.AssetDTO, int64, error)
}

type assetService struct {
	assetRepository repositories.AssetRepository
}

func NewAssetService(assetRepository repositories.AssetRepository) AssetService {
	return &assetService{assetRepository: assetRepository}
}

func (s *assetService) Save(assetDTO *dtos.AssetDTO) error {
	if !assetDTO.IsUserInputValid() {
		return fmt.Errorf("Invalid asset data: %v", assetDTO.ToMap())
	}

	asset := dtoToModel(assetDTO)
	if err := s.assetRepository.Save(asset); err != nil {
		return err
	}
	assetDTO.ID = asset.ID
	return nil
}

func (s *assetService) FindByCodeAndUserID(code string, userID uint) (*dtos.AssetDTO, error) {
	asset, err := s.assetRepository.FindByCodeAndUserID(code, userID)
	if err != nil {
		return nil, fmt.Errorf("Asset not found: %w", err)
	}
	return modelToDTO(asset), nil
}

func (s *assetService) GetPaginatedByUserID(
	userID uint, page, pageSize int,
) ([]*dtos.AssetDTO, int64, error) {
	assets, total, err := s.assetRepository.GetPaginatedByUserID(userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to get paginated assets: %w", err)
	}

	assetDTOs := make([]*dtos.AssetDTO, 0, len(assets))
	for i, asset := range assets {
		assetDTOs[i] = modelToDTO(asset)
	}
	return assetDTOs, total, nil
}

func modelToDTO(asset *models.Asset) *dtos.AssetDTO {
	if asset == nil {
		return nil
	}
	return &dtos.AssetDTO{
		ID:                 asset.ID,
		Code:               asset.Code,
		AssetType:          asset.AssetType,
		Currency:           asset.Currency,
		UserID:             asset.UserID,
		CachedHoldAvgPrice: asset.CachedHoldAvgPrice,
		CachedHoldQuantity: asset.CachedHoldQuantity,
	}
}

func dtoToModel(assetDTO *dtos.AssetDTO) *models.Asset {
	if assetDTO == nil {
		return nil
	}
	asset := &models.Asset{
		Code:               assetDTO.Code,
		AssetType:          assetDTO.AssetType,
		Currency:           assetDTO.Currency,
		UserID:             assetDTO.UserID,
		CachedHoldAvgPrice: assetDTO.CachedHoldAvgPrice,
		CachedHoldQuantity: assetDTO.CachedHoldQuantity,
	}
	asset.ID = assetDTO.ID
	return asset
}
