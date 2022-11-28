package main

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/mehmetokdemir/currency-conversion-service/config"
	_ "github.com/mehmetokdemir/currency-conversion-service/docs"
	"github.com/mehmetokdemir/currency-conversion-service/internal/account"
	"github.com/mehmetokdemir/currency-conversion-service/internal/currency"
	"github.com/mehmetokdemir/currency-conversion-service/internal/exchange"
	"github.com/mehmetokdemir/currency-conversion-service/internal/user"
	"github.com/mehmetokdemir/currency-conversion-service/middleware"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
	"log"
	"time"

	// External imports
	"github.com/patrickmn/go-cache"
)

// @title Currency Conversion Service
// @version 1.0.8
// @description Currency Conversion Service.
// @BasePath /
func main() {

	// Load Config
	serviceConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Currency Service
	currencyService := currency.NewCurrencyService(cache.New(5*time.Minute, 10*time.Minute))
	currencyService.SetLocalCacheToCurrencies()

	db := config.Connect(serviceConfig)

	// Account Service
	accountRepository := account.NewAccountRepository(db)
	if err = accountRepository.Migration(); err != nil {
		log.Fatal(err)
	}
	accountService := account.NewAccountService(accountRepository, serviceConfig)
	accountHandler := account.NewAccountHandler(accountService, currencyService)

	// User Service
	userRepository := user.NewUserRepository(db)
	if err = userRepository.Migration(); err != nil {
		log.Fatal(err)
	}
	userService := user.NewUserService(userRepository, serviceConfig, currencyService, accountService)
	userHandler := user.NewUserHandler(userService)

	// Exchange Service
	exchangeRepository := exchange.NewExchangeRepository(db)
	if err = exchangeRepository.Migration(); err != nil {
		log.Fatal(err)
	}
	exchangeService := exchange.NewExchangeService(exchangeRepository, currencyService, accountService)
	exchangeHandler := exchange.NewExchangeHandler(currencyService, exchangeService)

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
