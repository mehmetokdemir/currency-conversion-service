package main

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/mehmetokdemir/currency-conversion-service/config"
	_ "github.com/mehmetokdemir/currency-conversion-service/docs"
	"github.com/mehmetokdemir/currency-conversion-service/internal/handlers"
	"github.com/mehmetokdemir/currency-conversion-service/internal/repositories"
	"github.com/mehmetokdemir/currency-conversion-service/internal/services"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
	"log"
	"time"

	// External imports
	"github.com/patrickmn/go-cache"
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

	currencyService := services.NewCurrencyService(cache.New(5*time.Minute, 10*time.Minute))
	currencyService.SetLocalCacheToCurrencies()

	userRepository := repositories.NewUserRepository(config.Connect(serviceConfig.Db))
	if err = userRepository.Migration(); err != nil {
		log.Fatal(err)
	}

	userService := services.NewUserService(userRepository, serviceConfig)
	userHandler := handlers.NewUserHandler(userService, currencyService)

	// ACCOUNT SERVICE

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
