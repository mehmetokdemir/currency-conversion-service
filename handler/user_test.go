package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mehmetokdemir/currency-conversion-service/config"
	"github.com/mehmetokdemir/currency-conversion-service/entity"
	"github.com/mehmetokdemir/currency-conversion-service/helper"
	"github.com/mehmetokdemir/currency-conversion-service/repository"
	"github.com/mehmetokdemir/currency-conversion-service/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserHandler_Register(t *testing.T) {
	config, err := config.LoadConfig()
	assert.Nil(t, err)

	db := config.Db.Connect()
	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo, config)
	handler := NewUserHandler(service)
	app := gin.New()
	app.POST("/user/register", handler.Register)
	user := entity.User{
		Username: "test",
		Email:    "test@gmail.com",
	}
	user.Password, err = helper.HashPassword(user.Password)
	assert.Nil(t, err)

	created, err := repo.Create(user)
	assert.Nil(t, err)

	assert.Equal(t, "test", created.Username)

}
