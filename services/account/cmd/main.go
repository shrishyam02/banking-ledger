package main

import (
	"account/api"
	"account/processor"
	"account/repository"
	"account/service"
	"context"
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/shrishyam02/banking-ledger/common/config"
	"github.com/shrishyam02/banking-ledger/common/db"
	"github.com/shrishyam02/banking-ledger/common/kafka"
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

	brokers := strings.Split(cfg.Kafka.Brokers, ",")
	consumerTopics := []string{"account-balance-updates-topic"}
	consumerGroup := "account-group"
	producerTopics := []string{"transactions-status-topic"}

	for key, topic := range consumerTopics {
		logger.Log.Info().Msgf("Key: %d Topic: %s kafka broker %s", key, topic, brokers[0])
		kerr := kafka.CreateKafkaTopic(brokers[0], topic)
		if kerr != nil {
			logger.Log.Fatal().AnErr("Failed to create Kafka topic", kerr).Msg("test")
		}
		logger.Log.Info().Msgf("kafka broker string slice %s", brokers)
	}

	for key, topic := range producerTopics {
		logger.Log.Info().Msgf("Key: %d Topic: %s kafka broker %s", key, topic, brokers[0])
		kerr := kafka.CreateKafkaTopic(brokers[0], topic)
		if kerr != nil {
			logger.Log.Fatal().AnErr("Failed to create Kafka topic", kerr).Msg("test")
		}
		logger.Log.Info().Msgf("kafka broker string slice %s", brokers)
	}

	consumer := kafka.NewKafkaConsumer(brokers, consumerGroup, consumerTopics)
	producer := kafka.NewKafkaProducer(brokers)

	pgDb, err := db.ConnectPostgres(cfg.Database.PostgresConnectionString)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	logger.Log.Info().Msg("Connected to postgres: " + config.AccountService)

	accountRepo := repository.NewAccountRepository(pgDb)
	accountService := service.NewService(accountRepo)
	accountHandler := api.NewAccountHandler(accountService)
	accProcessor := processor.NewProcessor(consumer, producer, consumerTopics, producerTopics, consumerGroup, accountService)

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
	go func() {
		logger.Log.Info().Msg("Starting to process Account Balance")
		if err := accProcessor.ProcessAccountBalanceUpdates(ctx); err != nil {
			log.Fatalf("Failed to process account balance updates: %v", err)
		}
	}()
	server.RunServer(ctx, serverConfig, registerHandlers)
}
