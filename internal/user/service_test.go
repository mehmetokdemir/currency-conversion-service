package user

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/mehmetokdemir/currency-conversion-service/config"
	"github.com/mehmetokdemir/currency-conversion-service/internal/account"
	"github.com/mehmetokdemir/currency-conversion-service/internal/currency"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
	uService := NewUserService(mockUserRepository, config.Config{}, currency.Service{}, accService)

	request := User{
		Id:                  uint(1),
		Username:            "test",
		Email:               "test@gmail.com",
		Password:            "123",
		DefaultCurrencyCode: "TRY",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	expectedCreateUserResponse := &User{
		Id:       request.Id,
		Username: request.Username,
		Email:    request.Email,
		//Password:            "$2a$10$b86VMJDERpYuxCvVf6JvPuBHR7ipUQZzXmNm0/S3iA8Jp8votYY5y",
		Password:            "$2a$10$zX0LxdfVFt8brzbXe2bPOeOkXwncb3PcLxdZOTX6cxi5KE9tBlgwG",
		DefaultCurrencyCode: request.DefaultCurrencyCode,
		CreatedAt:           request.CreatedAt,
		UpdatedAt:           request.UpdatedAt,
	}

	mockUserRepository.EXPECT().IsUserExistWithSameUsername(request.Username).Return(false)
	mockUserRepository.EXPECT().IsUserExistWithSameEmail(request.Email).Return(false)
	mockUserRepository.EXPECT().CreateUser(request).Return(expectedCreateUserResponse, nil)

	user, err := uService.CreateUser(request)
	fmt.Println("user", user)
	assert.Nil(t, err)
	assert.Equal(t, expectedCreateUserResponse.Id, user.Id)
}
