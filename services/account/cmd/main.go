package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/shrishyam02/banking-ledger/common/config"
	"github.com/shrishyam02/banking-ledger/common/db"
	"github.com/shrishyam02/banking-ledger/common/logger"
)

func main() {
	logger.InitLogger()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	pgDb, err := db.ConnectPostgres(cfg.Database.PostgresConnectionString)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	// err = pgDb.AutoMigrate(&account.Account{})
	// if err != nil {
	// 	logger.Log.Fatal().Err(err).Msg("Failed to migrate database schema")
	// }

	// accountRepo := repository.NewAccountRepository(pgDb)
	// accountService := service.NewService(accountRepo)
	accountHandler := api.NewAccountHandler(accountService)

	router := gin.Default()
	apiGroup := router.Group("/api/v1")
	{
		accounts := apiGroup.Group("/accounts")
		{
			accounts.POST("", accountHandler.CreateAccount)
			accounts.GET("", accountHandler.ListAccounts)
			accounts.GET("/:id", accountHandler.GetAccount)
			// accounts.PUT("/:id", accountHandler.UpdateAccount)
			// accounts.DELETE("/:id", accountHandler.DeleteAccount)
			// accounts.GET("/number/:accountNumber", accountHandler.GetAccountByAccountNumber)
		}
	}

	server := &http.Server{
		Addr:    ":" + cfg.Services[config.AccountService].Port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("Server shutdown failed")
	}

	logger.Log.Println("Server gracefully stopped.")
}
