package account

import (
	"github.com/golang/mock/gomock"
	"github.com/mehmetokdemir/currency-conversion-service/config"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAccountService_CreateUserAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAccountRepository := NewMockIAccountRepository(ctrl)
	accService := NewAccountService(mockAccountRepository, config.Config{})
	request := Account{
		CurrencyCode: "TRY",
		UserId:       uint(1),
		Balance:      float64(300),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	expectedCreateAccountResponse := &Account{
		CurrencyCode: request.CurrencyCode,
		UserId:       request.UserId,
		Balance:      request.Balance,
		CreatedAt:    request.CreatedAt,
		UpdatedAt:    request.UpdatedAt,
	}

	t.Run("On registration balance set 10000", func(t *testing.T) {
		mockAccountRepository.EXPECT().CreateAccount(request).Return(expectedCreateAccountResponse, nil)
		account, err := accService.CreateUserAccount(request.UserId, request.CurrencyCode, true)
		assert.Nil(t, err)
		assert.Equal(t, expectedCreateAccountResponse.Balance, account.Balance)
		//assert.Less(t, expectedCreateAccountResponse.Balance, account.Balance)
	})

	t.Run("On try to conversion and not account on currency; account created and balance set 0", func(t *testing.T) {
		mockAccountRepository.EXPECT().CreateAccount(request).Return(expectedCreateAccountResponse, nil)
		account, err := accService.CreateUserAccount(request.UserId, request.CurrencyCode, false)
		assert.Nil(t, err)
		assert.NotEqual(t, expectedCreateAccountResponse.Balance, account.Balance)
		assert.Greater(t, expectedCreateAccountResponse.Balance, account.Balance)
	})

}

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

}
