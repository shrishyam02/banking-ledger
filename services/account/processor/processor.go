package processor

import (
	"context"
	"encoding/json"
	"sync"

	"account/service"

	ckafka "github.com/shrishyam02/banking-ledger/common/kafka"
	"github.com/shrishyam02/banking-ledger/common/logger"

	"github.com/segmentio/kafka-go"
)

type Processor interface {
	ProcessAccountBalanceUpdates(ctx context.Context) error
}

type processor struct {
	consumer       ckafka.KafkaConsumer
	producer       ckafka.KafkaProducer
	consumerTopics []string
	producerTopics []string
	consumerGroup  string
	accountService service.AccountService
}

func NewProcessor(consumer ckafka.KafkaConsumer, producer ckafka.KafkaProducer, consumerTopics []string, producerTopics []string, consumerGroup string, accountService service.AccountService) Processor {
	return &processor{
		consumer:       consumer,
		producer:       producer,
		consumerTopics: consumerTopics,
		producerTopics: producerTopics,
		consumerGroup:  consumerGroup,
		accountService: accountService,
	}
}

func (p *processor) ProcessAccountBalanceUpdates(ctx context.Context) error {
	logger.Log.Info().Msg("ProcessAccountBalanceUpdates")
	msgChan := make(chan kafka.Message)
	var wg sync.WaitGroup

	// Start worker pool
	for i := 0; i < 5; i++ { // 5 workers
		wg.Add(1)
		go p.worker(ctx, msgChan, &wg)
	}

	// Consume messages and send them to the worker pool
	go func() {
		defer close(msgChan)
		if err := p.consumer.Consume(ctx, p.consumerTopics[0], p.consumerGroup, func(msg kafka.Message) error {
			msgChan <- msg
			return nil
		}); err != nil {
			logger.Log.Fatal().Err(err).Msg("Failed to consume messages")
		}
	}()

	wg.Wait()
	return nil
}

func (p *processor) worker(ctx context.Context, msgChan <-chan kafka.Message, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range msgChan {
		if err := p.handleAccountBalanceUpdate(ctx, msg); err != nil {
			logger.Log.Error().Msgf("Failed to handle message: %v", err)
		}
	}
}

func (p *processor) handleAccountBalanceUpdate(ctx context.Context, msg kafka.Message) error {
	logger.Log.Info().Msg("handleAccountBalanceUpdate")

	var transaction map[string]interface{}
	if err := json.Unmarshal(msg.Value, &transaction); err != nil {
		return err
	}

	accountID := transaction["accountId"].(string)
	amount := transaction["amount"].(float64)
	transactionType := transaction["transactionType"].(string)
	logger.Log.Info().Msgf("handleAccountBalanceUpdate account balance pre update. account:(%v %v %v)", accountID, amount, transactionType)

	//TODO: retry logic on optimistic lock failure error
	err := p.accountService.UpdateAccountBalance(ctx, accountID, amount, transactionType)
	if err != nil {
		transaction["status"] = "failed"
		transaction["error"] = err.Error()
		logger.Log.Error().Msgf("handleAccountBalanceUpdate error while updating account balance account:(%v) err:%v", transaction, err.Error())
	} else {
		transaction["status"] = "success"
		transaction["error"] = ""
		logger.Log.Info().Msgf("handleAccountBalanceUpdate account balance update sucesss. account:(%v)", transaction)
	}

	transactionBytes, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	statusMessage := kafka.Message{
		Key:   msg.Key,
		Value: transactionBytes,
	}
	return p.producer.Produce(ctx, p.producerTopics[0], statusMessage)
}
