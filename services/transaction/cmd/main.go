package main

import (
	"context"
	"log"
	"os"
	"strings"

	"transaction/api"
	"transaction/service"

	"github.com/gin-gonic/gin"

	"github.com/shrishyam02/banking-ledger/common/config"
	ckafka "github.com/shrishyam02/banking-ledger/common/kafka"
	"github.com/shrishyam02/banking-ledger/common/logger"
	"github.com/shrishyam02/banking-ledger/common/server"
)

func main() {
	logger.InitLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	logger.Log.Info().Msg("Intialized logger for service: " + config.TransactionService)

	// Create Kafka topic
	brokers := strings.Split(cfg.Kafka.Brokers, ",")
	producerTopics := []string{"transactions-topic"}
	logger.Log.Info().Msgf("kafka broker %s", brokers[0])
	for key, topic := range producerTopics {
		logger.Log.Info().Msgf("Key: %d Topic: %s kafka broker %s", key, topic, brokers[0])
		kerr := ckafka.CreateKafkaTopic(brokers[0], topic)
		if kerr != nil {
			logger.Log.Fatal().AnErr("Failed to create Kafka topic", kerr).Msg("test")
		}
		logger.Log.Info().Msgf("kafka broker string slice %s", brokers)
	}
	producer := ckafka.NewKafkaProducer(brokers)

	accountServiceURL := os.Getenv("ACCOUNT_SERVICE_URL")
	if accountServiceURL == "" {
		log.Fatal("ACCOUNT_SERVICE_URL environment variable is required")
	}
	accountService := service.NewAccountService(accountServiceURL)
	transactionHandler := api.NewTransactionHandler(producer, producerTopics, accountService)

	registerHandlers := func(apiGroup *gin.RouterGroup) {
		accounts := apiGroup.Group("/transactions")
		{
			accounts.POST("", transactionHandler.CreateTransaction)
		}
	}
	logger.Log.Info().Msg("Handlers for: " + config.TransactionService)

	serverConfig := server.Config{
		Port:        cfg.Services[config.TransactionService].Port,
		ServiceName: config.TransactionService,
		ApiAuth:     cfg.ApiAuth,
	}

	ctx := context.Background()
	server.RunServer(ctx, serverConfig, registerHandlers)
}
