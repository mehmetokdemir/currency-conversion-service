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
	"time"
)

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	UserRoutes(router *gin.RouterGroup)
}

type userHandler struct {
	userService     services.UserService
	currencyService services.CurrencyService
}

func NewUserHandler(userService services.UserService, currencyService services.CurrencyService) UserHandler {
	return &userHandler{userService: userService, currencyService: currencyService}
}

func (h *userHandler) UserRoutes(router *gin.RouterGroup) {
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)
}

// Register godoc
// @Summary Create User
// @Description Create a user
// @Param request body dto.RegisterRequest true "body params"
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Response{data=dto.RegisterResponse} "Success"
// @Failure 400
// @Failure 500
// @Router /user/register [post]
func (h *userHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
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

	createdUser, err := h.userService.CreateUser(entity.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
	})
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, errors.ErrCreateError.Error(), err.Error())
		return
	}

	// TODO: CREATE WALLET ACCOUNT
	helper.Success(c, dto.RegisterResponse{Username: createdUser.Username, Email: createdUser.Email})
	return
}

// Login godoc
// @Summary Auth User
// @Description User Login
// @Param request body dto.LoginRequest true "body params"
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Response
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /user/login [post]
func (h *userHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
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

	rsp, err := h.userService.CreateToken(req.Username, req.Password)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, errors.ErrCreateTokenError.Error(), err.Error())
		return
	}

	helper.Success(c, rsp)
	return
}
