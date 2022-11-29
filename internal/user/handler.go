package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/mehmetokdemir/currency-conversion-service/errors"
	"github.com/mehmetokdemir/currency-conversion-service/helper"
	//"github.com/mehmetokdemir/currency-conversion-service/internal/account"
	"net/http"
	"time"
)

type Handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	UserRoutes(router *gin.RouterGroup)
}

type userHandler struct {
	userService IUserService
}

func NewUserHandler(userService IUserService) Handler {
	return &userHandler{userService: userService}
}

func (h *userHandler) UserRoutes(router *gin.RouterGroup) {
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)
}

// Register godoc
// @Summary Create User
// @Description Create a user
// @Tags User
// @Accept  json
// @Produce  json
// @Param request body RegisterRequest true "body params"
// @Success 200 {object} helper.Response{data=RegisterResponse} "Success"
// @Failure 400 {object} helper.Response{error=helper.ResponseError} "Bad Request"
// @Failure 404 {object} helper.Response{error=helper.ResponseError} "Not Found"
// @Failure 500 {object} helper.Response{error=helper.ResponseError} "Internal Server Error"
// @Router /user/register [post]
func (h *userHandler) Register(c *gin.Context) {
	var req RegisterRequest
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

	createdUser, err := h.userService.CreateUser(User{
		Username:            req.Username,
		Email:               req.Email,
		Password:            req.Password,
		DefaultCurrencyCode: req.CurrencyCode,
		CreatedAt:           time.Now().Local(),
		UpdatedAt:           time.Now().Local(),
	})
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, errors.ErrCreateError.Error(), err.Error())
		return
	}

	helper.Success(c, RegisterResponse{Username: createdUser.Username, Email: createdUser.Email})
}

// Login godoc
// @Summary Auth User
// @Description User Login
// @Tags User
// @Accept  json
// @Produce  json
// @Param request body LoginRequest true "body params"
// @Success 200 {object} helper.Response{data=LoginResponse} "Success"
// @Failure 400 {object} helper.Response{error=helper.ResponseError} "Bad Request"
// @Failure 404 {object} helper.Response{error=helper.ResponseError} "Not Found"
// @Failure 500 {object} helper.Response{error=helper.ResponseError} "Internal Server Error"
// @Router /user/login [post]
func (h *userHandler) Login(c *gin.Context) {
	var req LoginRequest
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
		helper.Error(c, http.StatusNotFound, errors.ErrCreateTokenError.Error(), err.Error())
		return
	}

	helper.Success(c, rsp)
}
