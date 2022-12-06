package user

import (
	// Go imports
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	// External imports
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserService := NewMockIUserService(ctrl)
	handler := NewUserHandler(mockUserService)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/register", handler.Register)
	t.Run("Status bad request", func(t *testing.T) {
		register := RegisterRequest{
			Username:     "test",
			Email:        "",
			Password:     "",
			CurrencyCode: "",
		}
		reqBytes, _ := json.Marshal(register)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBytes))
		if err != nil {
			t.Fatalf("Could not create request: %v\n", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Status OK", func(t *testing.T) {
		register := RegisterRequest{
			Username:     "test",
			Email:        "test@email.com",
			Password:     "123",
			CurrencyCode: "TRY",
		}

		expected := User{
			Id:                  uint(1),
			Username:            register.Username,
			Email:               register.Email,
			Password:            register.Password,
			DefaultCurrencyCode: register.CurrencyCode,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		mockUserService.EXPECT().CreateUser(expected).Return(&expected, nil)

		reqBytes, _ := json.Marshal(register)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBytes))
		if err != nil {
			t.Fatalf("Could not create request: %v\n", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		fmt.Println("w", w.Body)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestUserHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserService := NewMockIUserService(ctrl)
	httpHandler := NewUserHandler(mockUserService)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/login", httpHandler.Login)

	t.Run("Status bad request", func(t *testing.T) {
		login := LoginRequest{
			Username: "test",
		}

		reqBytes, _ := json.Marshal(login)
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBytes))
		if err != nil {
			t.Fatalf("Could not create request: %v\n", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Status not found", func(t *testing.T) {
		login := LoginRequest{
			Username: "test",
			Password: "test",
		}

		mockUserService.EXPECT().CreateToken(login.Username, login.Password).Return(nil, errors.New("user not found"))

		reqBytes, _ := json.Marshal(login)
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBytes))
		if err != nil {
			t.Fatalf("Could not create request: %v\n", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Status OK", func(t *testing.T) {
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
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBytes))
		if err != nil {
			t.Fatalf("Could not create request: %v\n", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, login.Username, expected.Username)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
