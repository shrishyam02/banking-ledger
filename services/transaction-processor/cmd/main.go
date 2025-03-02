package main

import (
	"context"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/shrishyam02/banking-ledger/common/config"
	"github.com/shrishyam02/banking-ledger/common/kafka"
	"github.com/shrishyam02/banking-ledger/common/logger"
	"github.com/shrishyam02/banking-ledger/common/server"

	"transaction-processor/processor"
)

func main() {
	logger.InitLogger()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	logger.Log.Info().Msg("Intialized logger for service: " + config.ProcessorService)

	// Create Kafka topic
	brokers := strings.Split(cfg.Kafka.Brokers, ",")
	consumerTopics := map[string]string{
		"transactions":        "transactions-topic",        //consumer
		"transactions-status": "transactions-status-topic", //consumer
	}
	producerTopics := map[string]string{
		"account-balance-updates": "account-balance-updates-topic", //producer
		"ledger":                  "ledger-topic",                  //producer
	}
	consumerGroup := "transaction-processor-group"
	for key, topic := range consumerTopics {
		logger.Log.Info().Msgf("Key: %s Topic: %s kafka broker %s", key, topic, brokers[0])
		kerr := kafka.CreateKafkaTopic(brokers[0], topic)
		if kerr != nil {
			logger.Log.Fatal().AnErr("Failed to create Kafka topic", kerr).Msg("test")
		}
		logger.Log.Info().Msgf("kafka broker string slice %s", brokers)
	}

	for key, topic := range producerTopics {
		logger.Log.Info().Msgf("Key: %s Topic: %s kafka broker %s", key, topic, brokers[0])
		kerr := kafka.CreateKafkaTopic(brokers[0], topic)
		if kerr != nil {
			logger.Log.Fatal().AnErr("Failed to create Kafka topic", kerr).Msg("test")
		}
		logger.Log.Info().Msgf("kafka broker string slice %s", brokers)
	}

	consumerTopicList := make([]string, 0, len(consumerTopics))
	for _, topic := range consumerTopics {
		consumerTopicList = append(consumerTopicList, topic)
	}
	consumer := kafka.NewKafkaConsumer(brokers, consumerGroup, consumerTopicList)
	producer := kafka.NewKafkaProducer(brokers)

	transactionProcessor := processor.NewTransactionProcessor(consumer, producer, consumerTopics, producerTopics, consumerGroup, 5) // 5 workers

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := transactionProcessor.ProcessTransactions(ctx); err != nil {
			log.Fatalf("Failed to process transactions: %v", err)
		}
	}()

	go func() {
		if err := transactionProcessor.ProcessTransactionStatus(ctx); err != nil {
			log.Fatalf("Failed to process transaction status: %v", err)
		}
	}()

	serverConfig := server.Config{
		Port:        cfg.Services[config.ProcessorService].Port,
		ServiceName: config.ProcessorService,
		ApiAuth:     cfg.ApiAuth,
	}

	server.RunServer(ctx, serverConfig, nil)
}
