package user

import (
	// Go imports
	"errors"
	"testing"
	"time"

	// External imports
	"github.com/golang/mock/gomock"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"

	// Internal imports
	"github.com/mehmetokdemir/currency-conversion-service/config"
	"github.com/mehmetokdemir/currency-conversion-service/internal/account"
	"github.com/mehmetokdemir/currency-conversion-service/internal/currency"
)

func TestUserService_CreateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserRepository := NewMockIUserRepository(ctrl)
	accService := account.NewMockIAccountService(ctrl)
	uService := NewUserService(mockUserRepository, config.Config{}, currency.Service{}, accService)

	t.Run("user not found error", func(t *testing.T) {
		username := "test-1"
		password := ""

		mockUserRepository.EXPECT().
			GetUserByUsername(username).
			Return(nil, errors.New("user_not_found"))

		authToken, err := uService.CreateToken(username, password)
		assert.NotNil(t, err)
		assert.Empty(t, authToken)
	})

	t.Run("password miss match error", func(t *testing.T) {
		username := "test-1"
		password := "123"

		expectedUserResp := &User{
			Id:        1,
			Username:  "test-1",
			Password:  "12345",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserRepository.EXPECT().
			GetUserByUsername(username).
			Return(expectedUserResp, errors.New("password_miss_match_error"))

		authToken, err := uService.CreateToken(username, password)
		assert.NotNil(t, err)
		assert.Empty(t, authToken)
	})

	t.Run("create token successfully", func(t *testing.T) {
		username := "test-1"
		password := "123"

		expectedUserResp := &User{
			Id:        1,
			Username:  "test-1",
			Password:  "$2a$10$mItUkS/0sjKfOCKCvJoDR.YrHcyv9.8/obPP0a91sPrF2aE7kjC76",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserRepository.EXPECT().
			GetUserByUsername(username).
			Return(expectedUserResp, nil)

		authToken, err := uService.CreateToken(username, password)
		assert.Nil(t, err)
		assert.NotEmpty(t, authToken)
	})

}

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserRepository := NewMockIUserRepository(ctrl)
	accService := account.NewMockIAccountService(ctrl)
	currencyService := currency.NewCurrencyService(cache.New(5*time.Minute, 10*time.Minute))
	currencyService.SetLocalCacheToCurrencies()
	mockUService := NewMockIUserService(ctrl)
	uService := NewUserService(mockUserRepository, config.Config{}, currencyService, accService)

	t.Run("Duplicated user error with same email", func(t *testing.T) {
		request := User{
			Id:                  uint(1),
			Username:            "test",
			Email:               "test@gmail.com",
			Password:            "123",
			DefaultCurrencyCode: "TRY",
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		mockUserRepository.EXPECT().IsUserExistWithSameEmail(request.Email).Return(true)
		_, err := uService.CreateUser(request)
		assert.NotNil(t, err)
		assert.Error(t, err, "duplicated users")
	})

	t.Run("Duplicated user error with same username", func(t *testing.T) {
		request := User{
			Id:                  uint(1),
			Username:            "test",
			Email:               "test@gmail.com",
			Password:            "123",
			DefaultCurrencyCode: "TRY",
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		mockUserRepository.EXPECT().IsUserExistWithSameEmail(request.Email).Return(false)
		mockUserRepository.EXPECT().IsUserExistWithSameUsername(request.Username).Return(true)
		_, err := uService.CreateUser(request)
		assert.NotNil(t, err)
		assert.Error(t, err, "duplicated users")
	})

	t.Run("Invalid currency code error", func(t *testing.T) {
		request := User{
			Id:                  uint(1),
			Username:            "test",
			Email:               "test@gmail.com",
			Password:            "123",
			DefaultCurrencyCode: "AGH",
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		mockUserRepository.EXPECT().IsUserExistWithSameEmail(request.Email).Return(false)
		mockUserRepository.EXPECT().IsUserExistWithSameUsername(request.Username).Return(false)
		ok := currencyService.CheckIsCurrencyCodeExist(request.DefaultCurrencyCode)
		assert.False(t, ok)
		_, err := uService.CreateUser(request)
		assert.NotNil(t, err)
		assert.Error(t, err, "currency not found")
	})

	t.Run("Successfully created user", func(t *testing.T) {
		request := User{
			Id:                  uint(1),
			Username:            "test",
			Email:               "test@gmail.com",
			Password:            "123",
			DefaultCurrencyCode: "TRY",
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		mockUserRepository.EXPECT().IsUserExistWithSameEmail(request.Email).Return(false)
		mockUserRepository.EXPECT().IsUserExistWithSameUsername(request.Username).Return(false)
		ok := currencyService.CheckIsCurrencyCodeExist(request.DefaultCurrencyCode)
		assert.True(t, ok)

		hashedPassword, err := uService.HashPassword(request.Password)
		assert.Nil(t, err)

		expectedUser := User{
			Id:                  request.Id,
			Username:            request.Username,
			Email:               request.Email,
			Password:            request.Password,
			DefaultCurrencyCode: request.DefaultCurrencyCode,
			CreatedAt:           request.CreatedAt,
			UpdatedAt:           request.UpdatedAt,
		}

		mockUserRepository.EXPECT().CreateUser(expectedUser).Return(&User{}, nil)
		mockUService.EXPECT().HashPassword(request.Password).Return(hashedPassword, nil)
		//accService.EXPECT().CreateUserAccount(expectedUser.Id, expectedUser.DefaultCurrencyCode, true).Return(nil, nil)
		user, err := uService.CreateUser(request)
		if err != nil {
			t.Errorf("Not error expected, but got: %q", err.Error())
		}
		assert.Equal(t, user.Id, expectedUser.Id)
	})
}

func TestUserService_MatchPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserRepository := NewMockIUserRepository(ctrl)
	accService := account.NewMockIAccountService(ctrl)
	currencyService := currency.NewCurrencyService(cache.New(5*time.Minute, 10*time.Minute))
	currencyService.SetLocalCacheToCurrencies()
	uService := NewUserService(mockUserRepository, config.Config{}, currencyService, accService)

	pass := "123"
	t.Run("match password", func(t *testing.T) {
		hashedPass, err := uService.HashPassword(pass)
		assert.Nil(t, err)
		ok := uService.VerifyPassword(hashedPass, pass)
		assert.True(t, ok)
	})

	t.Run("miss match password", func(t *testing.T) {
		hashedPass, err := uService.HashPassword(pass)
		assert.Nil(t, err)
		ok := uService.VerifyPassword(hashedPass, "1234")
		assert.False(t, ok)
	})
}
