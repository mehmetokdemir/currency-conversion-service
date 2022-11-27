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
	"github.com/mehmetokdemir/currency-conversion-service/middleware"
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

	// Load Config
	serviceConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal("12", err)
	}

	// Currency Service
	currencyService := services.NewCurrencyService(cache.New(5*time.Minute, 10*time.Minute))
	currencyService.SetLocalCacheToCurrencies()

	db := config.Connect(serviceConfig)

	// Account Service
	accountRepository := repositories.NewAccountRepository(db)
	if err = accountRepository.Migration(); err != nil {
		log.Fatal(err)
	}
	accountService := services.NewAccountService(accountRepository, serviceConfig)
	accountHandler := handlers.NewAccountHandler(accountService, currencyService)

	// User Service
	userRepository := repositories.NewUserRepository(db)
	if err = userRepository.Migration(); err != nil {
		log.Fatal(err)
	}
	userService := services.NewUserService(userRepository, serviceConfig, currencyService)
	userHandler := handlers.NewUserHandler(userService, accountService)

	// Exchange Service
	exchangeRepository := repositories.NewExchangeRepository(db)
	if err = exchangeRepository.Migration(); err != nil {
		log.Fatal(err)
	}
	exchangeService := services.NewExchangeService(exchangeRepository, currencyService, accountService)
	exchangeHandler := handlers.NewExchangeHandler(currencyService, exchangeService)

	// Gin App
	router := gin.New()
	router.Use(gin.Recovery())

	// User Routes
	userGroup := router.Group("/user")
	userGroup.Use()
	{
		userHandler.UserRoutes(userGroup)
	}

	// Account Routes
	accountGroup := router.Group("/account")
	accountGroup.Use(middleware.AuthMiddleware())
	{
		accountHandler.AccountRoutes(accountGroup)
	}

	// Exchange Routes
	exchangeGroup := router.Group("/exchange")
	exchangeGroup.Use(middleware.AuthMiddleware())
	{
		exchangeHandler.ExchangeRoutes(exchangeGroup)
	}

	// Swagger Documentation
	router.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))

	if err = router.Run(fmt.Sprintf(":%s", serviceConfig.ServerPort)); err != nil {
		log.Fatal(err)
	}
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}
