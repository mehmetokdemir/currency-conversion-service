package handlers

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/mehmetokdemir/currency-conversion-service/dto"
	"github.com/mehmetokdemir/currency-conversion-service/entity"
	"github.com/mehmetokdemir/currency-conversion-service/errors"
	"github.com/mehmetokdemir/currency-conversion-service/helper"
	"github.com/mehmetokdemir/currency-conversion-service/internal/services"
	"net/http"
)

type ExchangeHandler interface {
	ExchangeRate(c *gin.Context)
	ExchangeRoutes(router *gin.RouterGroup)
}

type exchangeHandler struct {
	currencyService services.CurrencyService
	exchangeService services.ExchangeService
}

func NewExchangeHandler(currencyService services.CurrencyService, exchangeService services.ExchangeService) ExchangeHandler {
	return &exchangeHandler{currencyService: currencyService, exchangeService: exchangeService}
}

func (h *exchangeHandler) ExchangeRoutes(router *gin.RouterGroup) {
	router.POST("/rate", h.ExchangeRate)
}

// ExchangeRate godoc
// @Summary Get Exchange Rate
// @Description Get exchange rate on given currencies
// @Tags Exchange
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param request body dto.ExchangeRateOfferRequest true "body params"
// @Success 200 {object} helper.Response{data=dto.ExchangeRateOfferResponse} "Success"
// @Failure 400 {object} helper.Response{error=helper.ResponseError} "Bad Request"
// @Failure 403 {object} helper.Response{error=helper.ResponseError} "Forbidden"
// @Failure 404 {object} helper.Response{error=helper.ResponseError} "Not Found"
// @Failure 500 {object} helper.Response{error=helper.ResponseError} "Internal Server Error"
// @Router /exchange/rate [post]
func (h *exchangeHandler) ExchangeRate(c *gin.Context) {
	var req dto.ExchangeRateOfferRequest
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

	userInContext, ok := c.Get("user")
	if !ok || userInContext == nil {
		helper.Error(c, http.StatusNotFound, errors.ErrNotFoundError.Error(), "user not found in context")
		return
	}

	user, ok := userInContext.(entity.User)
	if !ok {
		helper.Error(c, http.StatusBadRequest, errors.ErrDataTypeError.Error(), "type assertion error")
		return
	}

	exchangeRateResponse, err := h.exchangeService.GetExchangeRateOffer(user.ID, req)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, errors.ErrExchangeOfferError.Error(), err.Error())
		return
	}

	helper.Success(c, exchangeRateResponse)
}
