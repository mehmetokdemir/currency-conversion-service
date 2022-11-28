package account

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/mehmetokdemir/currency-conversion-service/internal/currency"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAccountService := NewMockIAccountService(ctrl)
	httpHandler := NewAccountHandler(mockAccountService, currency.Service{})
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/list", httpHandler.List)

	t.Run("user not found in context", func(t *testing.T) {
		reqTest := httptest.NewRequest(http.MethodPost, "/list", nil)
		reqTest.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodPost, "/list", nil)
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Auth-Token", "Invalid Token")
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("success", func(t *testing.T) {

	})
}
