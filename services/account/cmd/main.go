package main

import (
	"account/api"
	"account/repository"
	"account/service"
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/shrishyam02/banking-ledger/common/config"
	"github.com/shrishyam02/banking-ledger/common/db"
	"github.com/shrishyam02/banking-ledger/common/logger"
	"github.com/shrishyam02/banking-ledger/common/server"
)

func main() {
	logger.InitLogger()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	logger.Log.Info().Msg("Intialized logger for service: " + config.AccountService)

	pgDb, err := db.ConnectPostgres(cfg.Database.PostgresConnectionString)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	logger.Log.Info().Msg("Connected to postgres: " + config.AccountService)

	accountRepo := repository.NewAccountRepository(pgDb)
	accountService := service.NewService(accountRepo)
	accountHandler := api.NewAccountHandler(accountService)

	registerHandlers := func(apiGroup *gin.RouterGroup) {
		accounts := apiGroup.Group("/accounts")
		{
			accounts.POST("", accountHandler.CreateAccount)
			accounts.GET("", accountHandler.ListAccounts)
			accounts.GET("/:id", accountHandler.GetAccount)
		}
	}
	logger.Log.Info().Msg("Handlers for: " + config.AccountService)

	serverConfig := server.Config{
		Port:        cfg.Services[config.AccountService].Port,
		ServiceName: config.AccountService,
		ApiAuth:     cfg.ApiAuth,
	}

	ctx := context.Background()
	server.RunServer(ctx, serverConfig, registerHandlers)
}
