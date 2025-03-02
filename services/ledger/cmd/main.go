package main

import (
	"context"
	"ledger/api"
	"ledger/service"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/shrishyam02/banking-ledger/common/config"
	"github.com/shrishyam02/banking-ledger/common/db"
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
	logger.Log.Info().Msg("Initialized logger for service: " + config.LedgerService)

	brokers := strings.Split(cfg.Kafka.Brokers, ",")
	consumerTopics := []string{"ledger-topic"}
	consumerGroup := "ledger-group"

	for key, topic := range consumerTopics {
		logger.Log.Info().Msgf("Key: %d Topic: %s kafka broker %s", key, topic, brokers[0])
		kerr := ckafka.CreateKafkaTopic(brokers[0], topic)
		if kerr != nil {
			logger.Log.Fatal().AnErr("Failed to create Kafka topic", kerr).Msg("test")
		}
		logger.Log.Info().Msgf("kafka broker string slice %s", brokers)
	}

	consumer := ckafka.NewKafkaConsumer(brokers, consumerGroup, consumerTopics)

	mongoClient, err := db.ConnectMongo(cfg.Database.MongoDBConnectionString)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}
	logger.Log.Info().Msg("Connected to MongoDB: " + config.LedgerService)

	mongoDB := mongoClient.Database("ledger")
	ledgerService := service.NewLedgerService(mongoDB)
	ledgerHandler := api.NewLedgerHandler(ledgerService)

	registerHandlers := func(apiGroup *gin.RouterGroup) {
		ledger := apiGroup.Group("/ledger")
		{
			ledger.GET("/transactions/:id", ledgerHandler.GetTransactionHistory)
		}
	}
	logger.Log.Info().Msg("Handlers for: " + config.LedgerService)

	serverConfig := server.Config{
		Port:        cfg.Services[config.LedgerService].Port,
		ServiceName: config.LedgerService,
		ApiAuth:     cfg.ApiAuth,
	}

	ctx := context.Background()

	go func() {
		logger.Log.Info().Msg("Starting to process ledger transactions")
		if err := consumer.Consume(ctx, consumerTopics[0], consumerGroup, func(msg kafka.Message) error {
			return ledgerService.HandleMessage(ctx, msg)
		}); err != nil {
			log.Fatalf("Failed to process ledger transactions: %v", err)
		}
	}()
	server.RunServer(ctx, serverConfig, registerHandlers)
}
