package account

import (
	"github.com/gin-gonic/gin"
	"github.com/mehmetokdemir/currency-conversion-service/errors"
	"github.com/mehmetokdemir/currency-conversion-service/helper"
	"github.com/mehmetokdemir/currency-conversion-service/internal/common"
	"github.com/mehmetokdemir/currency-conversion-service/internal/currency"
	"net/http"
)

type Handler interface {
	List(c *gin.Context)
	AccountRoutes(router *gin.RouterGroup)
}

type accountHandler struct {
	accountService  IAccountService
	currencyService currency.Service
}

func NewAccountHandler(accountService IAccountService, currencyService currency.Service) Handler {
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
// @Success 200 {object} helper.Response{data=[]WalletAccount} "Success"
// @Failure 400 {object} helper.Response{error=helper.ResponseError} "Bad Request"
// @Failure 403 {object} helper.Response{error=helper.ResponseError} "Forbidden"
// @Failure 404 {object} helper.Response{error=helper.ResponseError} "Not Found"
// @Failure 500 {object} helper.Response{error=helper.ResponseError} "Internal Server Error"
// @Router /account/list [get]
func (h *accountHandler) List(c *gin.Context) {
	userId, ok := common.GetUserIdFromContext(c)
	if !ok {
		helper.Error(c, http.StatusNotFound, errors.ErrNotFoundError.Error(), "can not get user from context")
		return
	}

	walletAccounts, err := h.accountService.ListUserAccounts(userId)
	if err != nil {
		helper.Error(c, http.StatusNotFound, errors.ErrNotFoundError.Error(), err.Error())
		return
	}

	helper.Success(c, walletAccounts)
}
