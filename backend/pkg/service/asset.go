package service

import (
	"easyinvesting/pkg/dto"
	"easyinvesting/pkg/model"
	"easyinvesting/pkg/repository"
	"fmt"
)

type AssetService interface {
	Save(assetDTO *dto.AssetDTO) error
	FindByCodeAndUserID(code string, userID uint) (*dto.AssetDTO, error)
	GetPaginatedByUserID(userID uint, page, pageSize int) ([]*dto.AssetDTO, int64, error)
	DoesUserOwnAsset(code string, userID uint) (bool, error)
}

type assetService struct {
	assetRepository repository.AssetRepository
}

func NewAssetService(assetRepository repository.AssetRepository) AssetService {
	return &assetService{assetRepository: assetRepository}
}

func (s *assetService) Save(assetDTO *dto.AssetDTO) error {
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

func (s *assetService) FindByCodeAndUserID(code string, userID uint) (*dto.AssetDTO, error) {
	asset, err := s.assetRepository.FindByCodeAndUserID(code, userID)
	if err != nil {
		return nil, fmt.Errorf("Asset not found: %w", err)
	}
	return modelToDTO(asset), nil
}

func (s *assetService) GetPaginatedByUserID(
	userID uint, page, pageSize int,
) ([]*dto.AssetDTO, int64, error) {
	assets, total, err := s.assetRepository.GetPaginatedByUserID(userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to get paginated assets: %w", err)
	}

	assetDTOs := make([]*dto.AssetDTO, len(assets))
	for i, asset := range assets {
		assetDTOs[i] = modelToDTO(asset)
	}
	return assetDTOs, total, nil
}

func modelToDTO(asset *model.Asset) *dto.AssetDTO {
	if asset == nil {
		return nil
	}
	return &dto.AssetDTO{
		ID:                 asset.ID,
		Code:               asset.Code,
		AssetType:          asset.AssetType,
		Currency:           asset.Currency,
		UserID:             asset.UserID,
		CachedHoldAvgPrice: asset.CachedHoldAvgPrice,
		CachedHoldQuantity: asset.CachedHoldQuantity,
	}
}

func dtoToModel(assetDTO *dto.AssetDTO) *model.Asset {
	if assetDTO == nil {
		return nil
	}
	asset := &model.Asset{
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

func (s *assetService) DoesUserOwnAsset(code string, userID uint) (bool, error) {
	asset, err := s.assetRepository.FindByCodeAndUserID(code, userID)
	if err != nil {
		if err.Error() == "Asset not found" {
			return false, nil // User does not own the asset
		}
		return false, fmt.Errorf("Failed to check asset ownership: %w", err)
	}
	return asset != nil, nil // User owns the asset if it exists
}
