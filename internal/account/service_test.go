package account

import (
	// Go imports
	"errors"
	"testing"

	// External imports
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	// Internal imports
	"github.com/mehmetokdemir/currency-conversion-service/config"
)

func TestAccountService_GetUserBalanceOnGivenCurrencyAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAccountRepository := NewMockIAccountRepository(ctrl)
	accService := NewAccountService(mockAccountRepository, config.Config{})

	userId := uint(1)
	currencyCode := "USD"
	balance := float64(50)
	mockAccountRepository.EXPECT().GetUserBalanceOnGivenCurrencyAccount(userId, currencyCode).Return(balance, nil)

	actualBalance, err := accService.GetUserBalanceOnGivenCurrencyAccount(userId, currencyCode)
	assert.Nil(t, err)
	assert.Equal(t, actualBalance, balance)
}

func TestAccountService_IsUserHasAccountOnGivenCurrency(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAccountRepository := NewMockIAccountRepository(ctrl)
	accService := NewAccountService(mockAccountRepository, config.Config{})
	userId := uint(1)
	currencyCode := "USD"
	t.Run("user has account on given currency", func(t *testing.T) {
		mockAccountRepository.EXPECT().IsUserHasAccountOnGivenCurrency(userId, currencyCode).Return(true)
		ok := accService.IsUserHasAccountOnGivenCurrency(userId, currencyCode)
		assert.True(t, ok)
	})

	t.Run("user has not account on given currency", func(t *testing.T) {
		mockAccountRepository.EXPECT().IsUserHasAccountOnGivenCurrency(userId, "USD").Return(false)
		ok := accService.IsUserHasAccountOnGivenCurrency(userId, currencyCode)
		assert.False(t, ok)
	})
}

func TestAccountService_UpdateUserBalanceOnGivenCurrencyAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAccountRepository := NewMockIAccountRepository(ctrl)
	accService := NewAccountService(mockAccountRepository, config.Config{})

	userId := uint(1)
	currencyCode := "USD"
	balance := float64(50)

	t.Run("account not found on given currency", func(t *testing.T) {
		mockAccountRepository.EXPECT().GetUserBalanceOnGivenCurrencyAccount(userId, currencyCode).Return(float64(0), errors.New("account not found"))
		err := accService.UpdateUserBalanceOnGivenCurrencyAccount(userId, currencyCode, balance)
		assert.NotNil(t, err)
	})

	t.Run("successfully updated balance", func(t *testing.T) {
		mockAccountRepository.EXPECT().GetUserBalanceOnGivenCurrencyAccount(userId, currencyCode).Return(balance, nil)
		mockAccountRepository.EXPECT().UpdateUserBalanceOnGivenCurrencyAccount(userId, currencyCode, balance).Return(nil)
		balance -= balance
		err := accService.UpdateUserBalanceOnGivenCurrencyAccount(userId, currencyCode, balance)
		assert.Nil(t, err)
	})
}
