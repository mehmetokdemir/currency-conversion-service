package dto

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/mehmetokdemir/currency-conversion-service/entity"
)

type Token struct {
	User entity.User
	jwt.StandardClaims
}
