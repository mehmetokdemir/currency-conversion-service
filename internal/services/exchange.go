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
	AcceptExchangeRateOffer(userId uint, request dto.ExchangeAcceptOfferRequest) ([]dto.AccountWallet, error)
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
	offerId, err := s.createExchangeRateOffer(userId, fromCurrencyCode, toCurrencyCode, exchangeRateWithMarkupRate)
	if err != nil {
		return nil, err
	}

	return &dto.ExchangeRateOfferResponse{
		OfferId:          offerId,
		FromCurrencyCode: exchange.FromCurrencyCode,
		ToCurrencyCode:   exchange.ToCurrencyCode,
		ExchangeRate:     exchangeRateWithMarkupRate,
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

func (s *exchangeService) AcceptExchangeRateOffer(userId uint, request dto.ExchangeAcceptOfferRequest) ([]dto.AccountWallet, error) {
	offer, err := s.exchangeRepo.GetOffer(request.OfferId)
	if err != nil {
		return nil, err
	}

	// Check offer for user
	if offer.UserId != userId {
		return nil, errors.New("exchange rate is not valid for this user")
	}

	// Check offer is valid
	if offer.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("offer has expired")
	}

	balance, err := s.accountService.GetUserBalanceOnGivenCurrencyAccount(userId, offer.FromCurrencyCode)
	if err != nil {
		return nil, err
	}

	if request.Amount > balance {
		return nil, errors.New("not enough balance")
	}

	if err = s.updateUserBalances(userId, *offer, request.Amount); err != nil {
		return nil, err
	}

	accountsWithBalances, err := s.accountService.ListUserAccounts(userId)
	if err != nil {
		return nil, err
	}

	return accountsWithBalances, nil
}

func (s *exchangeService) updateUserBalances(userId uint, offer entity.Offer, amount float64) error {
	// Calculate from currency balance
	fromCurrencyCode, fromBalance := s.calculateFromBalanceAfterAcceptedCurrencyConversion(offer, amount)
	if err := s.accountService.UpdateUserBalanceOnGivenCurrencyAccount(userId, fromCurrencyCode, fromBalance); err != nil {
		return err
	}

	// Calculate to currency balance
	toCurrencyCode, toBalance := s.calculateToBalanceAfterAcceptedCurrencyConversion(offer, amount)
	if err := s.accountService.UpdateUserBalanceOnGivenCurrencyAccount(userId, toCurrencyCode, toBalance); err != nil {
		return err
	}

	return nil
}

func (s *exchangeService) calculateToBalanceAfterAcceptedCurrencyConversion(offer entity.Offer, amount float64) (string, float64) {
	return offer.ToCurrencyCode, amount * offer.ExchangeRate
}

func (s *exchangeService) calculateFromBalanceAfterAcceptedCurrencyConversion(offer entity.Offer, amount float64) (string, float64) {
	return offer.FromCurrencyCode, -1 * amount
}
