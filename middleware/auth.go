package middleware

import (
	// Go imports
	"net/http"
	"os"

	// External imports
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	// Internal imports
	"github.com/mehmetokdemir/currency-conversion-service/dto"
	"github.com/mehmetokdemir/currency-conversion-service/errors"
	"github.com/mehmetokdemir/currency-conversion-service/helper"
)

func middlewareError(statusCode int, message string, detail string) helper.Response {
	return helper.Response{
		Success:    false,
		StatusCode: statusCode,
		Error: &helper.ResponseError{
			Message: message,
			Detail:  detail,
		},
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.GetHeader("X-Auth-Token") == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, middlewareError(http.StatusForbidden, errors.ErrNotFoundError.Error(), "token not found"))
			return
		}

		tk := dto.Token{}
		token, err := jwt.ParseWithClaims(c.GetHeader("X-Auth-Token"), &tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, middlewareError(http.StatusForbidden, errors.ErrExpiredTokenError.Error(), err.Error()))
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusForbidden, middlewareError(http.StatusForbidden, errors.ErrInvalidTokenError.Error(), "token is not valid"))
			return
		}

		c.Set("user_id", tk.UserId)

		c.Next()
	}
}
