package exchange

import (
	// Go imports
	"errors"
	"fmt"
	"testing"
	"time"

	// External imports
	"github.com/golang/mock/gomock"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"

	// Internal imports
	"github.com/mehmetokdemir/currency-conversion-service/internal/account"
	"github.com/mehmetokdemir/currency-conversion-service/internal/currency"
)

func TestExchangeService_AcceptExchangeRateOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockExchangeRepository := NewMockIExchangeRepository(ctrl)
	accService := account.NewMockIAccountService(ctrl)
	currencyService := currency.NewCurrencyService(cache.New(5*time.Minute, 10*time.Minute))
	currencyService.SetLocalCacheToCurrencies()
	exchService := NewExchangeService(mockExchangeRepository, currencyService, accService)

	t.Run("offer not found", func(t *testing.T) {
		userId := uint(1)
		acceptOfferRequest := AcceptOfferRequest{
			OfferId: uint(1),
			Amount:  100,
		}

		mockExchangeRepository.EXPECT().GetOffer(acceptOfferRequest.OfferId).Return(&Offer{}, errors.New("offer not found"))
		_, err := exchService.AcceptExchangeRateOffer(userId, acceptOfferRequest)
		assert.NotNil(t, err)
	})

	t.Run("offer expired", func(t *testing.T) {
		userId := uint(1)
		acceptOfferRequest := AcceptOfferRequest{
			OfferId: uint(1),
			Amount:  100,
		}

		expectedOffer := Offer{
			Id:               acceptOfferRequest.OfferId,
			FromCurrencyCode: "TRY",
			ToCurrencyCode:   "USD",
			ExchangeRate:     18.50,
			ExpiresAt:        time.Now().Add(time.Minute * -3).Unix(),
			UserId:           userId,
		}

		mockExchangeRepository.EXPECT().GetOffer(acceptOfferRequest.OfferId).Return(&expectedOffer, nil)
		_, err := exchService.AcceptExchangeRateOffer(userId, acceptOfferRequest)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "offer has expired")
	})

	t.Run("has no balance for conversion", func(t *testing.T) {
		userId := uint(1)
		acceptOfferRequest := AcceptOfferRequest{
			OfferId: uint(1),
			Amount:  100,
		}

		expectedOffer := Offer{
			Id:               acceptOfferRequest.OfferId,
			FromCurrencyCode: "TRY",
			ToCurrencyCode:   "USD",
			ExchangeRate:     18.50,
			ExpiresAt:        time.Now().Add(time.Minute * 3).Unix(),
			UserId:           userId,
		}

		mockExchangeRepository.EXPECT().GetOffer(acceptOfferRequest.OfferId).Return(&expectedOffer, nil)
		accService.EXPECT().GetUserBalanceOnGivenCurrencyAccount(userId, expectedOffer.FromCurrencyCode).Return(float64(50), nil)
		_, err := exchService.AcceptExchangeRateOffer(userId, acceptOfferRequest)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "not enough balance")
	})

	t.Run("accept conversion offer", func(t *testing.T) {
		userId := uint(1)
		acceptOfferRequest := AcceptOfferRequest{
			OfferId: uint(1),
			Amount:  100,
		}

		expectedOffer := Offer{
			Id:               acceptOfferRequest.OfferId,
			FromCurrencyCode: "TRY",
			ToCurrencyCode:   "USD",
			ExchangeRate:     18.50,
			ExpiresAt:        time.Now().Add(time.Minute * 3).Unix(),
			UserId:           userId,
		}

		mockExchangeRepository.EXPECT().GetOffer(acceptOfferRequest.OfferId).Return(&expectedOffer, nil)
		accService.EXPECT().GetUserBalanceOnGivenCurrencyAccount(userId, expectedOffer.FromCurrencyCode).Return(float64(150), nil)

		amount := acceptOfferRequest.Amount * expectedOffer.ExchangeRate
		accService.EXPECT().UpdateUserBalanceOnGivenCurrencyAccount(userId, expectedOffer.FromCurrencyCode, -1*acceptOfferRequest.Amount).Return(nil)
		accService.EXPECT().UpdateUserBalanceOnGivenCurrencyAccount(userId, expectedOffer.ToCurrencyCode, amount).Return(nil)

		accService.EXPECT().ListUserAccounts(userId).Return([]account.WalletAccount{}, nil)
		_, err := exchService.AcceptExchangeRateOffer(userId, acceptOfferRequest)
		assert.Nil(t, err)
	})

}

func TestExchangeService_GetExchangeRateOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockExchangeRepository := NewMockIExchangeRepository(ctrl)
	accService := account.NewMockIAccountService(ctrl)
	currencyService := currency.NewCurrencyService(cache.New(5*time.Minute, 10*time.Minute))
	currencyService.SetLocalCacheToCurrencies()
	exchService := NewExchangeService(mockExchangeRepository, currencyService, accService)

	t.Run("currency is not exist", func(t *testing.T) {
		fromCurrencyCode := "AAA"
		toCurrencyCode := "USD"

		isExistFromCurrencyCode := currencyService.CheckIsCurrencyCodeExist(fromCurrencyCode)
		assert.False(t, isExistFromCurrencyCode)

		isExistToCurrencyCode := currencyService.CheckIsCurrencyCodeExist(toCurrencyCode)
		assert.True(t, isExistToCurrencyCode)
	})

	t.Run("user has no account on given from currency code", func(t *testing.T) {
		userId := uint(1)
		fromCurrencyCode := "USD"

		accService.EXPECT().IsUserHasAccountOnGivenCurrency(userId, fromCurrencyCode).Return(false)
		assert.Error(t, fmt.Errorf("%s account not found", fromCurrencyCode))
	})

	t.Run("exchange rate offer success", func(t *testing.T) {
		userId := uint(1)
		now := time.Now()
		offerRequest := OfferRequest{
			FromCurrencyCode: "TRY",
			ToCurrencyCode:   "USD",
		}

		expectedExchangeResponse := &Exchange{
			FromCurrencyCode: offerRequest.FromCurrencyCode,
			ToCurrencyCode:   offerRequest.ToCurrencyCode,
			ExchangeRate:     18,
			MarkupRate:       2,
			CreatedAt:        now,
			UpdatedAt:        now,
		}

		accService.EXPECT().IsUserHasAccountOnGivenCurrency(userId, offerRequest.FromCurrencyCode).Return(true)
		accService.EXPECT().IsUserHasAccountOnGivenCurrency(userId, offerRequest.ToCurrencyCode).Return(true)
		mockExchangeRepository.EXPECT().GetExchangeRate(offerRequest.FromCurrencyCode, offerRequest.ToCurrencyCode).Return(expectedExchangeResponse, nil)
		offer := Offer{
			Id:               uint(2),
			FromCurrencyCode: expectedExchangeResponse.FromCurrencyCode,
			ToCurrencyCode:   expectedExchangeResponse.ToCurrencyCode,
			ExchangeRate:     expectedExchangeResponse.ExchangeRate - expectedExchangeResponse.MarkupRate,
			ExpiresAt:        time.Now().Add(time.Minute * 3).Unix(),
			UserId:           userId,
			CreatedAt:        now,
			UpdatedAt:        now,
		}

		offerId, err := exchService.CreateExchangeRateOffer(userId, offer.FromCurrencyCode, offer.ToCurrencyCode, offer.ExchangeRate)
		assert.Nil(t, err)

		off, err := exchService.GetExchangeRateOffer(userId, offerRequest)
		assert.Nil(t, err)
		assert.Equal(t, off.OfferId, offerId)
	})
}

func TestExchangeService_CreateExchangeRateOffer(t *testing.T) {
	//TODO:
	ctrl := gomock.NewController(t)
	mockExchangeRepository := NewMockIExchangeRepository(ctrl)
	accService := account.NewMockIAccountService(ctrl)
	currencyService := currency.NewCurrencyService(cache.New(5*time.Minute, 10*time.Minute))
	currencyService.SetLocalCacheToCurrencies()
	exchService := NewExchangeService(mockExchangeRepository, currencyService, accService)

	offerId := uint(1)
	expiresAt := time.Now().Add(time.Minute * 3).Unix()
	offer := Offer{
		Id:               offerId,
		FromCurrencyCode: "TRY",
		ToCurrencyCode:   "USD",
		ExchangeRate:     18,
		ExpiresAt:        expiresAt,
		UserId:           uint(1),
	}

	expectedOffer := &Offer{
		Id:               offerId,
		FromCurrencyCode: "TRY",
		ToCurrencyCode:   "USD",
		ExchangeRate:     18,
		ExpiresAt:        expiresAt,
		UserId:           uint(1),
	}

	mockExchangeRepository.EXPECT().CreateOffer(gomock.Any()).Return(expectedOffer, nil)
	_, err := exchService.CreateExchangeRateOffer(offer.UserId, offer.FromCurrencyCode, offer.ToCurrencyCode, offer.ExchangeRate)
	assert.NotNil(t, err)
}
