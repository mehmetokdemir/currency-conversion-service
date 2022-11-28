package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserService := NewMockIUserService(ctrl)
	httpHandler := NewUserHandler(mockUserService)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/register", httpHandler.Register)
	now := time.Now()
	t.Run("Status bad request", func(t *testing.T) {
		register := RegisterRequest{
			Username:     "test",
			Email:        "",
			Password:     "",
			CurrencyCode: "",
		}

		reqBytes, _ := json.Marshal(register)
		reqTest := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBytes))
		reqTest.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBytes))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Success", func(t *testing.T) {
		actual := User{
			Id:                  1,
			Username:            "test",
			Email:               "test@gmail.com",
			Password:            "12345",
			DefaultCurrencyCode: "TRY",
			CreatedAt:           now,
			UpdatedAt:           now,
		}

		expected := &User{
			Id:                  1,
			Username:            "test",
			Email:               "test@gmail.com",
			Password:            "$2a$10$mItUkS/0sjKfOCKCvJoDR.YrHcyv9.8/obPP0a91sPrF2aE7kjC76",
			DefaultCurrencyCode: "TRY",
			CreatedAt:           now,
			UpdatedAt:           now,
			DeletedAt:           gorm.DeletedAt{},
		}

		mockUserService.EXPECT().CreateUser(actual).Return(expected, nil)

		reqBytes, _ := json.Marshal(RegisterRequest{
			Username:     actual.Username,
			Email:        actual.Email,
			Password:     actual.Password,
			CurrencyCode: actual.DefaultCurrencyCode,
		})
		reqTest := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBytes))
		reqTest.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBytes))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestUserHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserService := NewMockIUserService(ctrl)
	httpHandler := NewUserHandler(mockUserService)
	router := gin.Default()
	router.POST("/login", httpHandler.Login)
	gin.SetMode(gin.TestMode)
	t.Run("Status bad request", func(t *testing.T) {
		login := LoginRequest{
			Username: "test",
		}

		reqBytes, _ := json.Marshal(login)
		reqTest := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBytes))
		reqTest.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBytes))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Status not found", func(t *testing.T) {
		login := LoginRequest{
			Username: "test",
			Password: "test",
		}

		mockUserService.EXPECT().CreateToken(login.Username, login.Password).Return(&LoginResponse{
			RegisterResponse: RegisterResponse{
				Username: "test1",
				Email:    "test1",
			},
			TokenHash: "",
		}, nil)

		reqBytes, _ := json.Marshal(login)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBytes))
		assert.Nil(t, err)
		router.ServeHTTP(rr, req)
		fmt.Println("code", rr.Code)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Success", func(t *testing.T) {
		login := LoginRequest{
			Username: "test",
			Password: "pass",
		}

		expected := &LoginResponse{
			RegisterResponse: RegisterResponse{
				Username: "test",
				Email:    "pass",
			},
			TokenHash: "token",
		}

		mockUserService.EXPECT().CreateToken(login.Username, login.Password).Return(expected, nil)

		reqBytes, _ := json.Marshal(login)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBytes))

		assert.Nil(t, err)
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
