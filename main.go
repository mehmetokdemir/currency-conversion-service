package main

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/mehmetokdemir/currency-conversion-service/config"
	_ "github.com/mehmetokdemir/currency-conversion-service/docs"
	"github.com/mehmetokdemir/currency-conversion-service/handler"
	"github.com/mehmetokdemir/currency-conversion-service/repository"
	"github.com/mehmetokdemir/currency-conversion-service/service"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
	"log"
)

// @title Currency Conversion Service
// @version 1.0
// @description Currency Conversion Service.
// @BasePath /
func main() {

	serviceConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db := serviceConfig.Db.Connect()

	userRepository := repository.NewUserRepository(db)
	if err = userRepository.Migration(); err != nil {
		log.Fatal(err)
	}
	userService := service.NewUserService(userRepository, serviceConfig)
	userHandler := handler.NewUserHandler(userService)

	router := gin.New()
	router.Use(gin.Recovery())

	userGroup := router.Group("/user")
	userGroup.Use()
	{
		userHandler.UserRoutes(userGroup)
	}

	router.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))
	if err = router.Run(fmt.Sprintf(":%s", serviceConfig.Server.Port)); err != nil {
		log.Fatal(err)
	}
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}
