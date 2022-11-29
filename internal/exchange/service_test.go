package exchange

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/mehmetokdemir/currency-conversion-service/internal/account"
	"github.com/mehmetokdemir/currency-conversion-service/internal/currency"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExchangeService_AcceptExchangeRateOffer(t *testing.T) {
}

func TestExchangeService_GetExchangeRateOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockExchangeRepository := NewMockIExchangeRepository(ctrl)
	accService := account.NewMockIAccountService(ctrl)
	currencyService := currency.NewCurrencyService(cache.New(5*time.Minute, 10*time.Minute))
	currencyService.SetLocalCacheToCurrencies()

	// TODO: Make createExchangeRateOffer func as a public
	//exchService := NewExchangeService(mockExchangeRepository, currencyService, accService)

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
		fromCurrencyCode := "TRY"
		toCurrencyCode := "TRY"

		exchange := &Exchange{
			FromCurrencyCode: fromCurrencyCode,
			ToCurrencyCode:   toCurrencyCode,
			ExchangeRate:     18.63,
			MarkupRate:       1.0,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		expectedResponse := &OfferResponse{
			OfferId:          uint(1),
			FromCurrencyCode: fromCurrencyCode,
			ToCurrencyCode:   toCurrencyCode,
			ExchangeRate:     exchange.ExchangeRate - exchange.MarkupRate,
		}

		mockExchangeRepository.EXPECT().GetExchangeRate(toCurrencyCode, toCurrencyCode).Return(exchange, nil)
		assert.Equal(t, expectedResponse.ExchangeRate, 17.63)
	})

}
