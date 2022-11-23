package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mehmetokdemir/currency-conversion-service/dto"
	"github.com/mehmetokdemir/currency-conversion-service/errors"
	"github.com/mehmetokdemir/currency-conversion-service/helper"
	"net/http"
	"os"
	"strings"
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

		if c.GetHeader("X-Auth") == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, middlewareError(http.StatusForbidden, errors.ErrNotFoundError.Error(), "token not found"))
			return
		}

		tokenSplit := strings.Split(c.GetHeader("X-Auth"), " ")
		if len(tokenSplit) != 2 {
			c.AbortWithStatusJSON(http.StatusForbidden, middlewareError(http.StatusForbidden, errors.ErrInvalidJwtError.Error(), "jwt is invalid"))
			return
		}

		tokenPart := tokenSplit[1]
		tk := dto.Token{}
		token, err := jwt.ParseWithClaims(tokenPart, &tk, func(token *jwt.Token) (interface{}, error) {
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

		c.Set("user", tk.User)

		c.Next()
	}
}
