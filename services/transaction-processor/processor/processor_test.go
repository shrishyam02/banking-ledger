package processor

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKafkaConsumer struct {
	mock.Mock
}

func (m *MockKafkaConsumer) Consume(ctx context.Context, topic, groupID string, handler func(kafka.Message) error) error {
	args := m.Called(ctx, topic, groupID, handler)
	return args.Error(0)
}

type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) Produce(ctx context.Context, topic string, message kafka.Message) error {
	args := m.Called(ctx, topic, message)
	return args.Error(0)
}

func TestProcessTransactions(t *testing.T) {
	mockConsumer := new(MockKafkaConsumer)
	mockProducer := new(MockKafkaProducer)
	consumerTopics := map[string]string{"transactions": "transactions-topic"}
	producerTopics := map[string]string{"account-balance-updates": "account-balance-updates-topic", "ledger": "ledger-topic"}
	processor := NewTransactionProcessor(mockConsumer, mockProducer, consumerTopics, producerTopics, "consumer-group", 1)

	mockConsumer.On("Consume", mock.Anything, "transactions-topic", "consumer-group", mock.Anything).Return(nil)

	ctx := context.Background()
	err := processor.ProcessTransactions(ctx)
	assert.NoError(t, err)
	mockConsumer.AssertExpectations(t)
}

func TestProcessTransactionStatus(t *testing.T) {
	mockConsumer := new(MockKafkaConsumer)
	mockProducer := new(MockKafkaProducer)
	consumerTopics := map[string]string{"transactions-status": "transactions-status-topic"}
	producerTopics := map[string]string{"ledger": "ledger-topic"}
	processor := NewTransactionProcessor(mockConsumer, mockProducer, consumerTopics, producerTopics, "consumer-group", 1)

	mockConsumer.On("Consume", mock.Anything, "transactions-status-topic", "consumer-group", mock.Anything).Return(nil)

	ctx := context.Background()
	err := processor.ProcessTransactionStatus(ctx)
	assert.NoError(t, err)
	mockConsumer.AssertExpectations(t)
}

func TestHandleMessage(t *testing.T) {
	mockProducer := new(MockKafkaProducer)
	processor := &TransactionProcessor{
		producer:       mockProducer,
		producerTopics: map[string]string{"account-balance-updates": "account-balance-updates-topic", "ledger": "ledger-topic"},
	}

	transaction := map[string]interface{}{
		"amount": 100.0,
	}
	transactionBytes, _ := json.Marshal(transaction)
	msg := kafka.Message{
		Key:   []byte("key"),
		Value: transactionBytes,
	}

	mockProducer.On("Produce", mock.Anything, "account-balance-updates-topic", mock.Anything).Return(nil)

	ctx := context.Background()
	err := processor.handleMessage(ctx, msg)
	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestHandleMessage_InvalidTransaction(t *testing.T) {
	mockProducer := new(MockKafkaProducer)
	processor := &TransactionProcessor{
		producer:       mockProducer,
		producerTopics: map[string]string{"ledger": "ledger-topic"},
	}

	transaction := map[string]interface{}{
		"amount": -100.0,
	}
	transactionBytes, _ := json.Marshal(transaction)
	msg := kafka.Message{
		Key:   []byte("key"),
		Value: transactionBytes,
	}

	mockProducer.On("Produce", mock.Anything, "ledger-topic", mock.Anything).Return(nil)

	ctx := context.Background()
	err := processor.handleMessage(ctx, msg)
	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestHandleStatusMessage(t *testing.T) {
	mockProducer := new(MockKafkaProducer)
	processor := &TransactionProcessor{
		producer:       mockProducer,
		producerTopics: map[string]string{"ledger": "ledger-topic"},
	}

	transactionStatus := map[string]interface{}{
		"status": "completed",
	}
	transactionStatusBytes, _ := json.Marshal(transactionStatus)
	msg := kafka.Message{
		Key:   []byte("key"),
		Value: transactionStatusBytes,
	}

	mockProducer.On("Produce", mock.Anything, "ledger-topic", mock.Anything).Return(nil)

	ctx := context.Background()
	err := processor.handleStatusMessage(ctx, msg)
	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestValidateTransaction(t *testing.T) {
	processor := &TransactionProcessor{}

	validTransaction := map[string]interface{}{
		"amount": 100.0,
	}
	err := processor.validateTransaction(validTransaction)
	assert.NoError(t, err)

	invalidTransaction := map[string]interface{}{
		"amount": -100.0,
	}
	err = processor.validateTransaction(invalidTransaction)
	assert.Error(t, err)
	assert.Equal(t, "invalid transaction amount", err.Error())
}

func TestPublishTransactionStatus(t *testing.T) {
	mockProducer := new(MockKafkaProducer)
	processor := &TransactionProcessor{
		producer:       mockProducer,
		producerTopics: map[string]string{"ledger": "ledger-topic"},
	}

	transaction := map[string]interface{}{
		"amount": 100.0,
	}
	transactionBytes, _ := json.Marshal(transaction)
	msg := kafka.Message{
		Key:   []byte("key"),
		Value: transactionBytes,
	}

	mockProducer.On("Produce", mock.Anything, "ledger-topic", mock.Anything).Return(nil)

	ctx := context.Background()
	err := processor.publishTransactionStatus(ctx, msg.Key, transaction)
	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}
