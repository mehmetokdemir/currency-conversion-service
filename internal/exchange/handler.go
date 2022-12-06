package exchange

import (
	// Go imports
	"net/http"

	// External imports
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

	// Internal imports
	"github.com/mehmetokdemir/currency-conversion-service/errors"
	"github.com/mehmetokdemir/currency-conversion-service/helper"
	"github.com/mehmetokdemir/currency-conversion-service/internal/common"
	"github.com/mehmetokdemir/currency-conversion-service/internal/currency"
)

type Handler interface {
	ExchangeRate(c *gin.Context)
	AcceptOffer(c *gin.Context)
	ExchangeRoutes(router *gin.RouterGroup)
}

type exchangeHandler struct {
	currencyService currency.Service
	exchangeService IExchangeService
}

func NewExchangeHandler(currencyService currency.Service, exchangeService IExchangeService) Handler {
	return &exchangeHandler{currencyService: currencyService, exchangeService: exchangeService}
}

func (h *exchangeHandler) ExchangeRoutes(router *gin.RouterGroup) {
	router.POST("/rate", h.ExchangeRate)
	router.POST("/accept/offer", h.AcceptOffer)
}

// ExchangeRate godoc
// @Summary Get Exchange Rate
// @Description Get exchange rate on given currencies
// @Tags Exchange
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param request body OfferRequest true "body params"
// @Success 200 {object} helper.Response{data=OfferResponse} "Success"
// @Failure 400 {object} helper.Response{error=helper.ResponseError} "Bad Request"
// @Failure 403 {object} helper.Response{error=helper.ResponseError} "Forbidden"
// @Failure 404 {object} helper.Response{error=helper.ResponseError} "Not Found"
// @Failure 500 {object} helper.Response{error=helper.ResponseError} "Internal Server Error"
// @Router /exchange/rate [post]
func (h *exchangeHandler) ExchangeRate(c *gin.Context) {
	var req OfferRequest
	if err := c.BindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, errors.ErrBindJson.Error(), err.Error())
		return
	}

	_, err := govalidator.ValidateStruct(req)
	warnings := helper.WarningsFromValidationError(err)
	if warnings != nil {
		helper.Warning(c, warnings)
		return
	}

	userId, ok := common.GetUserIdFromContext(c)
	if !ok {
		helper.Error(c, http.StatusNotFound, errors.ErrNotFoundError.Error(), "can not get user from context")
		return
	}

	exchangeRateResponse, err := h.exchangeService.GetExchangeRateOffer(userId, req)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, errors.ErrExchangeOfferError.Error(), err.Error())
		return
	}

	helper.Success(c, exchangeRateResponse)
}

// AcceptOffer godoc
// @Summary Accept exchange rate offer
// @Description Accept the given exchange rate
// @Tags Exchange
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param request body AcceptOfferRequest true "body params"
// @Success 200 {object} helper.Response{data=[]account.WalletAccount} "Success"
// @Failure 400 {object} helper.Response{error=helper.ResponseError} "Bad Request"
// @Failure 403 {object} helper.Response{error=helper.ResponseError} "Forbidden"
// @Failure 404 {object} helper.Response{error=helper.ResponseError} "Not Found"
// @Failure 500 {object} helper.Response{error=helper.ResponseError} "Internal Server Error"
// @Router /exchange/accept/offer [post]
func (h *exchangeHandler) AcceptOffer(c *gin.Context) {
	var req AcceptOfferRequest
	if err := c.BindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, errors.ErrBindJson.Error(), err.Error())
		return
	}

	_, err := govalidator.ValidateStruct(req)
	warnings := helper.WarningsFromValidationError(err)
	if warnings != nil {
		helper.Warning(c, warnings)
		return
	}

	userId, ok := common.GetUserIdFromContext(c)
	if !ok {
		helper.Error(c, http.StatusNotFound, errors.ErrNotFoundError.Error(), "can not get user from context")
		return
	}

	accountsWithBalances, err := h.exchangeService.AcceptExchangeRateOffer(userId, req)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, errors.ErrExchangeOfferAcceptedError.Error(), err.Error())
		return
	}

	helper.Success(c, accountsWithBalances)
}
