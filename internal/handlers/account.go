package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mehmetokdemir/currency-conversion-service/dto"
	"github.com/mehmetokdemir/currency-conversion-service/errors"
	"github.com/mehmetokdemir/currency-conversion-service/helper"
	"github.com/mehmetokdemir/currency-conversion-service/internal/services"
	"net/http"
)

type AccountHandler interface {
	List(c *gin.Context)
	AccountRoutes(router *gin.RouterGroup)
}

type accountHandler struct {
	accountService  services.AccountService
	currencyService services.CurrencyService
}

func NewAccountHandler(accountService services.AccountService, currencyService services.CurrencyService) AccountHandler {
	return &accountHandler{accountService: accountService, currencyService: currencyService}
}

func (h *accountHandler) AccountRoutes(router *gin.RouterGroup) {
	router.GET("/list", h.List)
}

// List godoc
// @Summary List User Accounts
// @Description List user's balances with currencies
// @Tags Account
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Success 200 {object} helper.Response{data=dto.UserAccountResponse} "Success"
// @Failure 400 {object} helper.Response{error=helper.ResponseError} "Bad Request"
// @Failure 403 {object} helper.Response{error=helper.ResponseError} "Forbidden"
// @Failure 404 {object} helper.Response{error=helper.ResponseError} "Not Found"
// @Failure 500 {object} helper.Response{error=helper.ResponseError} "Internal Server Error"
// @Router /account/list [get]
func (h *accountHandler) List(c *gin.Context) {

	user, ok := getUserFromContext(c)
	if !ok {
		helper.Error(c, http.StatusNotFound, errors.ErrNotFoundError.Error(), "can not get user from context")
		return
	}

	walletAccounts, err := h.accountService.ListUserAccounts(user.ID)
	if err != nil {
		helper.Error(c, http.StatusNotFound, errors.ErrNotFoundError.Error(), err.Error())
		return
	}

	helper.Success(c, dto.UserAccountResponse{
		Id:       user.ID,
		Username: user.Username,
		Wallets:  walletAccounts,
	})
}
