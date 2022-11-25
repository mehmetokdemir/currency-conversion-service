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

	serviceConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// CURRENCY SERVICE
	currencyService := services.NewCurrencyService(cache.New(5*time.Minute, 10*time.Minute))
	currencyService.SetLocalCacheToCurrencies()

	// ACCOUNT SERVICE
	accountRepository := repositories.NewAccountRepository(config.Connect(serviceConfig.Db))
	if err = accountRepository.Migration(); err != nil {
		log.Fatal(err)
	}
	accountService := services.NewAccountService(accountRepository, serviceConfig)
	accountHandler := handlers.NewAccountHandler(accountService, currencyService)

	// USER SERVICE
	userRepository := repositories.NewUserRepository(config.Connect(serviceConfig.Db))
	if err = userRepository.Migration(); err != nil {
		log.Fatal(err)
	}
	userService := services.NewUserService(userRepository, serviceConfig, currencyService)
	userHandler := handlers.NewUserHandler(userService, accountService)

	// EXCHANGE SERVICE
	exchangeRepository := repositories.NewExchangeRepository(config.Connect(serviceConfig.Db))
	if err = exchangeRepository.Migration(); err != nil {
		log.Fatal(err)
	}
	exchangeService := services.NewExchangeService(exchangeRepository, currencyService, accountService)
	exchangeHandler := handlers.NewExchangeHandler(currencyService, exchangeService)

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

	// Exchange routes
	exchangeGroup := router.Group("/exchange")
	exchangeGroup.Use(middleware.AuthMiddleware())
	{
		exchangeHandler.ExchangeRoutes(exchangeGroup)
	}

	router.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))
	if err = router.Run(fmt.Sprintf(":%s", serviceConfig.Server.Port)); err != nil {
		log.Fatal(err)
	}
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}
