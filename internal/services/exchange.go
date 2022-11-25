package services

import (
	"errors"
	"fmt"
	"github.com/mehmetokdemir/currency-conversion-service/dto"
	"github.com/mehmetokdemir/currency-conversion-service/entity"
	"github.com/mehmetokdemir/currency-conversion-service/internal/repositories"
	"strings"
	"time"
)

type ExchangeService interface {
	GetExchangeRateOffer(userId uint, request dto.ExchangeRateOfferRequest) (*dto.ExchangeRateOfferResponse, error)
}

type exchangeService struct {
	exchangeRepo    repositories.ExchangeRepository
	currencyService CurrencyService
	accountService  AccountService
}

func NewExchangeService(exchangeRepository repositories.ExchangeRepository, currencyService CurrencyService, accountService AccountService) ExchangeService {
	return &exchangeService{exchangeRepo: exchangeRepository, currencyService: currencyService, accountService: accountService}
}

func (s *exchangeService) GetExchangeRateOffer(userId uint, request dto.ExchangeRateOfferRequest) (*dto.ExchangeRateOfferResponse, error) {
	fromCurrencyCode := strings.ToUpper(request.FromCurrencyCode)
	toCurrencyCode := strings.ToUpper(request.ToCurrencyCode)

	for _, currencyCode := range []string{fromCurrencyCode, toCurrencyCode} {
		if ok := s.currencyService.CheckCurrencyCodeIsExist(currencyCode); !ok {
			return nil, errors.New("currency not found")
		}
	}

	if ok := s.accountService.IsUserHasAccountOnGivenCurrency(userId, fromCurrencyCode); !ok {
		return nil, fmt.Errorf("%s account not found", fromCurrencyCode)
	}

	if ok := s.accountService.IsUserHasAccountOnGivenCurrency(userId, toCurrencyCode); !ok {
		if err := s.accountService.CreateUserAccount(userId, toCurrencyCode); err != nil {
			return nil, err
		}
	}

	exchange, err := s.exchangeRepo.GetExchangeRate(fromCurrencyCode, toCurrencyCode)
	if err != nil {
		return nil, err
	}

	exchangeRateWithMarkupRate := exchange.ExchangeRate + exchange.MarkupRate
	fmt.Println("exchange", exchange.ExchangeRate)
	fmt.Println("markup", exchange.MarkupRate)
	fmt.Println("tot", exchangeRateWithMarkupRate)
	offerId, err := s.createExchangeRateOffer(userId, fromCurrencyCode, toCurrencyCode, exchangeRateWithMarkupRate)
	if err != nil {
		return nil, err
	}

	return &dto.ExchangeRateOfferResponse{
		OfferId:          offerId,
		FromCurrencyCode: exchange.FromCurrencyCode,
		ToCurrencyCode:   exchange.ToCurrencyCode,
		//ExchangeRate:     math.Ceil((exchangeRateWithMarkupRate * 100) / 100),
		ExchangeRate: exchangeRateWithMarkupRate,
	}, nil
}

func (s *exchangeService) createExchangeRateOffer(userId uint, fromCurrencyCode, toCurrencyCode string, exchangeRate float64) (uint, error) {
	offer := entity.Offer{
		FromCurrencyCode: fromCurrencyCode,
		ToCurrencyCode:   toCurrencyCode,
		ExchangeRate:     exchangeRate,
		ExpiresAt:        time.Now().Add(time.Minute * 3).Unix(),
		UserId:           userId,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	offerId, err := s.exchangeRepo.CreateOffer(offer)
	if err != nil {
		return 0, err
	}

	return offerId, nil
}
