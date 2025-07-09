package service

import (
	"easyinvesting/pkg/client"
	"easyinvesting/pkg/dto"
	"easyinvesting/pkg/model"
	"easyinvesting/pkg/repository"
	"log"
	"net/http"
	"time"
)

type AssetOnMarketService interface {
	EnsureAssetOnMarket(code string) error
	UpdateAllAssetsOnMarket() error
	GetRealTimeAssetData(code string) (dto.DailyAssetPriceDTO, error)
}

type assetOnMarketService struct {
	assetOnMarketRepository   repository.AssetOnMarketRepository
	dailyAssetPriceRepository repository.DailyAssetPriceRepository
}

func NewAssetOnMarketService(
	assetOnMarketRepository repository.AssetOnMarketRepository,
	dailyAssetPriceRepository repository.DailyAssetPriceRepository,
) AssetOnMarketService {
	return &assetOnMarketService{
		assetOnMarketRepository:   assetOnMarketRepository,
		dailyAssetPriceRepository: dailyAssetPriceRepository,
	}
}

func (s *assetOnMarketService) EnsureAssetOnMarket(code string) error {
	created, err := s.assetOnMarketRepository.EnsureAssetOnMarket(code)
	if err != nil {
		return err
	}

	if created {
		// creates a new daily asset price entry for the asset
		brApiClient := client.NewBrApi(&http.Client{})
		quote, err := brApiClient.GetQuote(code)
		if err != nil && err != client.BrApiErrNoResults {
			return err
		}

		dailyAssetPrice := model.DailyAssetPrice{
			AssetCode: code,
			Price:     quote.RegularMarketPrice,
			Date:      time.Now().Format("2006-01-02"),
		}
		if err := s.dailyAssetPriceRepository.Save(&dailyAssetPrice); err != nil {
			log.Printf("Error creating daily asset price for %s: %v", code, err)
			return err
		}
	}

	return nil
}

func (s *assetOnMarketService) UpdateAllAssetsOnMarket() error {
	codes, err := s.assetOnMarketRepository.GetAllAssetCodes()
	if err != nil {
		return err
	}

	for _, code := range codes {
		if err := s.EnsureAssetOnMarket(code); err != nil {
			return err
		}
	}
	return s.dailyAssetPriceRepository.UpdateAllAssetsOnMarket(codes)
}

func (s *assetOnMarketService) GetRealTimeAssetData(code string) (dto.DailyAssetPriceDTO, error) {
	if err := s.EnsureAssetOnMarket(code); err != nil {
		log.Printf("Failed to ensure asset on market for code %s: %v", code, err)
		return dto.DailyAssetPriceDTO{}, err
	}

	dailyPrice, err := s.dailyAssetPriceRepository.FindLatestPriceByCode(code)
	if err != nil {
		return dto.DailyAssetPriceDTO{}, err
	}

	return dto.DailyAssetPriceDTO{
		AssetCode: dailyPrice.AssetCode,
		Price:     dailyPrice.Price,
		Date:      dailyPrice.Date,
	}, nil
}
