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
	userId := uint(1)
	auth := func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("user_id", userId)
			c.Next()
		}
	}

	t.Run("user not found in context", func(t *testing.T) {
		router.POST("/list", httpHandler.List)
		req, err := http.NewRequest(http.MethodPost, "/list", nil)
		if err != nil {
			t.Fatalf("Could not create request: %v\n", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("successfully list user wallet accounts", func(t *testing.T) {

		router.POST("/list", auth(httpHandler.List))
		walletAccounts := []WalletAccount{
			{
				CurrencyCode: "TRY",
				Balance:      50,
			},
			{
				CurrencyCode: "EUR",
				Balance:      35,
			},
			{
				CurrencyCode: "USD",
				Balance:      60,
			},
		}

		mockAccountService.EXPECT().ListUserAccounts(userId).Return(walletAccounts, nil)

		req, err := http.NewRequest(http.MethodPost, "/list", nil)
		if err != nil {
			t.Fatalf("Could not create request: %v\n", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Auth-Token", "app-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

	})
}
