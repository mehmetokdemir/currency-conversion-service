package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mehmetokdemir/currency-conversion-service/dto"
	"github.com/mehmetokdemir/currency-conversion-service/entity"
	"github.com/mehmetokdemir/currency-conversion-service/helper"
	"github.com/mehmetokdemir/currency-conversion-service/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Success", func(t *testing.T) {
		user := &entity.User{
			Username: "test",
			Email:    "test@gmail.com",
			Password: "1234",
		}

		mockRegisterResponse := helper.Response{
			Success:    true,
			StatusCode: http.StatusOK,
			Warnings:   nil,
			Error:      nil,
			Data: dto.RegisterResponse{
				Username: "test",
				Email:    "test@gmail.com",
			},
		}
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Register", mock.AnythingOfType("*dto.RegisterRequest"), user).Return(mockRegisterResponse, nil)
		rr := httptest.NewRecorder()
		router := gin.Default()

		pass, err := mockUserService.HashPassword(user.Password)
		assert.NoError(t, err)
		reqBody, err := json.Marshal(gin.H{
			"username": user.Username,
			"email":    user.Email,
			"password": pass,
		})
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, req)
		rsp, err := json.Marshal(mockRegisterResponse)
		fmt.Println("code", rr.Code)
		assert.NoError(t, err)
		assert.Equal(t, 200, rr.Code)
		assert.Equal(t, rsp, rr.Body.Bytes())
		mockUserService.AssertExpectations(t)
	})
}

func TestUserHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Success", func(t *testing.T) {
		user := &dto.LoginRequest{
			Username: "test",
			Password: "1234",
		}

		mockLoginResponse := helper.Response{
			Success:    true,
			StatusCode: http.StatusOK,
			Warnings:   nil,
			Error:      nil,
			Data: dto.LoginResponse{
				RegisterResponse: dto.RegisterResponse{
					Username: "test",
					Email:    "",
				},
				TokenHash: "",
			},
		}
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Login", mock.AnythingOfType("*dto.LoginRequest"), user).Return(mockLoginResponse, nil)
		rr := httptest.NewRecorder()
		router := gin.Default()

		pass, err := mockUserService.HashPassword(user.Password)
		assert.NoError(t, err)
		reqBody, err := json.Marshal(gin.H{
			"username": user.Username,
			"password": pass,
		})
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, req)
		rsp, err := json.Marshal(mockLoginResponse)
		fmt.Println("code", rr.Code)
		assert.NoError(t, err)
		assert.Equal(t, 200, rr.Code)
		assert.Equal(t, rsp, rr.Body.Bytes())
		mockUserService.AssertExpectations(t)
	})

}
