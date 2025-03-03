package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	ckafka "github.com/shrishyam02/banking-ledger/common/kafka"
)

type TransactionProcessor struct {
	consumer       ckafka.KafkaConsumer
	producer       ckafka.KafkaProducer
	consumerTopics map[string]string
	producerTopics map[string]string
	consumerGroup  string
	workerPoolSize int
}

func NewTransactionProcessor(consumer ckafka.KafkaConsumer, producer ckafka.KafkaProducer, consumerTopics map[string]string, producerTopics map[string]string, consumerGroup string, workerPoolSize int) *TransactionProcessor {
	return &TransactionProcessor{
		consumer:       consumer,
		producer:       producer,
		consumerTopics: consumerTopics,
		producerTopics: producerTopics,
		consumerGroup:  consumerGroup,
		workerPoolSize: workerPoolSize,
	}
}

func (tp *TransactionProcessor) ProcessTransactions(ctx context.Context) error {
	msgChan := make(chan kafka.Message)
	var wg sync.WaitGroup

	// Start worker pool
	for i := 0; i < tp.workerPoolSize; i++ {
		wg.Add(1)
		go tp.worker(ctx, msgChan, &wg)
	}

	// Consume messages and send them to the worker pool
	go func() {
		defer close(msgChan)
		if err := tp.consumer.Consume(ctx, tp.consumerTopics["transactions"], tp.consumerGroup, func(msg kafka.Message) error {
			msgChan <- msg
			return nil
		}); err != nil {
			log.Fatalf("Failed to consume messages: %v", err)
		}
	}()

	wg.Wait()
	return nil
}

func (tp *TransactionProcessor) ProcessTransactionStatus(ctx context.Context) error {
	msgChan := make(chan kafka.Message)
	var wg sync.WaitGroup

	// Start worker pool
	for i := 0; i < tp.workerPoolSize; i++ {
		wg.Add(1)
		go tp.statusWorker(ctx, msgChan, &wg)
	}

	// Consume messages and send them to the worker pool
	go func() {
		defer close(msgChan)
		if err := tp.consumer.Consume(ctx, tp.consumerTopics["transactions-status"], tp.consumerGroup, func(msg kafka.Message) error {
			msgChan <- msg
			return nil
		}); err != nil {
			log.Fatalf("Failed to consume messages: %v", err)
		}
	}()

	wg.Wait()
	return nil
}

func (tp *TransactionProcessor) worker(ctx context.Context, msgChan <-chan kafka.Message, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range msgChan {
		if err := tp.handleMessage(ctx, msg); err != nil {
			log.Printf("Failed to handle message: %v", err)
		}
	}
}

func (tp *TransactionProcessor) statusWorker(ctx context.Context, msgChan <-chan kafka.Message, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range msgChan {
		if err := tp.handleStatusMessage(ctx, msg); err != nil {
			log.Printf("Failed to handle status message: %v", err)
		}
	}
}

func (tp *TransactionProcessor) handleMessage(ctx context.Context, msg kafka.Message) error {
	var transaction map[string]interface{}
	if err := json.Unmarshal(msg.Value, &transaction); err != nil {
		return err
	}

	// Validate the transaction
	if err := tp.validateTransaction(transaction); err != nil {
		transaction["processedAt"] = time.Now().UTC()
		transaction["status"] = "failed"
		transaction["error"] = err.Error()
		return tp.publishTransactionStatus(ctx, msg.Key, transaction)
	}

	// Publish the transaction to account-service for balance update
	accountMessage := kafka.Message{
		Key:   msg.Key,
		Value: msg.Value,
	}
	if err := tp.producer.Produce(ctx, tp.producerTopics["account-balance-updates"], accountMessage); err != nil {
		transaction["processedAt"] = time.Now().UTC()
		transaction["status"] = "failed"
		transaction["error"] = err.Error()
		return tp.publishTransactionStatus(ctx, msg.Key, transaction)
	}

	return nil
}

func (tp *TransactionProcessor) handleStatusMessage(ctx context.Context, msg kafka.Message) error {
	var transactionStatus map[string]interface{}
	if err := json.Unmarshal(msg.Value, &transactionStatus); err != nil {
		return err
	}

	if err := tp.publishTransactionStatus(ctx, msg.Key, transactionStatus); err != nil {
		return err
	}

	return nil
}

func (tp *TransactionProcessor) validateTransaction(transaction map[string]interface{}) error {
	// TODO: Basic validation is added here. Additional validations need to be added.
	amount, ok := transaction["amount"].(float64)
	if !ok || amount <= 0 {
		return fmt.Errorf("invalid transaction amount")
	}
	return nil
}

func (tp *TransactionProcessor) publishTransactionStatus(ctx context.Context, key []byte, transaction map[string]interface{}) error {
	transactionBytes, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	ledgerMessage := kafka.Message{
		Key:   key,
		Value: transactionBytes,
	}
	return tp.producer.Produce(ctx, tp.producerTopics["ledger"], ledgerMessage)
}
