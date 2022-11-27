package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mehmetokdemir/currency-conversion-service/entity"
)

func getUserFromContext(c *gin.Context) (entity.User, bool) {
	userInContext, ok := c.Get("user")
	if !ok || userInContext == nil {
		return entity.User{}, false
	}

	user, ok := userInContext.(entity.User)
	if !ok {
		return entity.User{}, false
	}

	return user, true
}
