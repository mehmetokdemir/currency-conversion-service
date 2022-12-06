package common

import (
	// External imports
	"github.com/gin-gonic/gin"
)

func GetUserIdFromContext(c *gin.Context) (uint, bool) {
	userInContext, ok := c.Get("user_id")
	if !ok || userInContext == nil {
		return 0, false
	}

	userId, ok := userInContext.(uint)
	if !ok {
		return 0, false
	}

	return userId, true
}
