package exchange

import (
	// Go imports
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	// External imports
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	// Internal imports
	"github.com/mehmetokdemir/currency-conversion-service/internal/account"
	"github.com/mehmetokdemir/currency-conversion-service/internal/currency"
)

func TestExchangeHandler_AcceptOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockExchangeService := NewMockIExchangeService(ctrl)
	httpHandler := NewExchangeHandler(currency.Service{}, mockExchangeService)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	userId := uint(1)
	auth := func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("user_id", userId)
			c.Next()
		}
	}

	t.Run("user not in context", func(t *testing.T) {
		register := AcceptOfferRequest{
			OfferId: uint(1),
			Amount:  100,
		}
		router.POST("/accept/offer", httpHandler.AcceptOffer)
		reqBytes, _ := json.Marshal(register)
		req, err := http.NewRequest(http.MethodPost, "/accept/offer", bytes.NewReader(reqBytes))
		if err != nil {
			t.Fatalf("Could not create request: %v\n", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("invalid request", func(t *testing.T) {
		register := AcceptOfferRequest{
			OfferId: uint(1),
			Amount:  0,
		}
		router.POST("/accept/offer", httpHandler.AcceptOffer)
		reqBytes, _ := json.Marshal(register)

		mockExchangeService.EXPECT().AcceptExchangeRateOffer(userId, register).Return(nil, errors.New(""))

		req, err := http.NewRequest(http.MethodPost, "/accept/offer", bytes.NewReader(reqBytes))
		if err != nil {
			t.Fatalf("Could not create request: %v\n", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("offer successfully accepted", func(t *testing.T) {
		acceptOfferRequest := AcceptOfferRequest{
			OfferId: uint(1),
			Amount:  100,
		}
		router.POST("/accept/offer", auth(httpHandler.AcceptOffer))
		reqBytes, _ := json.Marshal(acceptOfferRequest)

		walletAccounts := []account.WalletAccount{
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

		mockExchangeService.EXPECT().AcceptExchangeRateOffer(userId, acceptOfferRequest).Return(walletAccounts, nil)

		req, err := http.NewRequest(http.MethodPost, "/accept/offer", bytes.NewReader(reqBytes))
		if err != nil {
			t.Fatalf("Could not create request: %v\n", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestExchangeHandler_ExchangeRate(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockExchangeService := NewMockIExchangeService(ctrl)
	httpHandler := NewExchangeHandler(currency.Service{}, mockExchangeService)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	userId := uint(1)
	auth := func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("user_id", userId)
			c.Next()
		}
	}

	t.Run("user not in context", func(t *testing.T) {
		register := OfferRequest{
			FromCurrencyCode: "TRY",
			ToCurrencyCode:   "USD",
		}
		router.POST("/exchange/rate", httpHandler.ExchangeRate)
		reqBytes, _ := json.Marshal(register)
		req, err := http.NewRequest(http.MethodPost, "/exchange/rate", bytes.NewReader(reqBytes))
		if err != nil {
			t.Fatalf("Could not create request: %v\n", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		fmt.Println("W", w)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		register := OfferRequest{
			FromCurrencyCode: "TRY",
		}

		router.POST("/exchange/rate", httpHandler.ExchangeRate)
		reqBytes, _ := json.Marshal(register)
		req, err := http.NewRequest(http.MethodPost, "/exchange/rate", bytes.NewReader(reqBytes))
		if err != nil {
			t.Fatalf("Could not create request: %v\n", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		fmt.Println("W", w)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("successfully get exchange rate for given currencies", func(t *testing.T) {
		offerRequest := OfferRequest{
			FromCurrencyCode: "TRY",
			ToCurrencyCode:   "USD",
		}

		router.POST("/exchange/rate", auth(httpHandler.ExchangeRate))
		reqBytes, _ := json.Marshal(offerRequest)

		expectedResponse := OfferResponse{
			OfferId:          uint(1),
			FromCurrencyCode: offerRequest.FromCurrencyCode,
			ToCurrencyCode:   offerRequest.ToCurrencyCode,
			ExchangeRate:     18.50,
		}

		mockExchangeService.EXPECT().GetExchangeRateOffer(userId, offerRequest).Return(&expectedResponse, nil)
		req, err := http.NewRequest(http.MethodPost, "/exchange/rate", bytes.NewReader(reqBytes))
		if err != nil {
			t.Fatalf("Could not create request: %v\n", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
