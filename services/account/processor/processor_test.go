package processor

import (
	"account/model"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKafkaConsumer struct {
	mock.Mock
}

func (m *MockKafkaConsumer) Consume(ctx context.Context, topic, group string, handler func(kafka.Message) error) error {
	args := m.Called(ctx, topic, group, handler)
	return args.Error(0)
}

type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) Produce(ctx context.Context, topic string, message kafka.Message) error {
	args := m.Called(ctx, topic, message)
	return args.Error(0)
}

type MockAccountService struct {
	mock.Mock
}

func (m *MockAccountService) UpdateAccountBalance(ctx context.Context, accountID string, amount float64, transactionType string) error {
	args := m.Called(ctx, accountID, amount, transactionType)
	return args.Error(0)
}

func (m *MockAccountService) CreateAccount(account *model.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountService) CreateOrUpdateCustomer(customer *model.Customer) error {
	args := m.Called(customer)
	return args.Error(0)
}

func (m *MockAccountService) GetAccountByID(accountID uuid.UUID) (*model.Account, error) {
	args := m.Called(accountID)
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockAccountService) ListAccounts() ([]model.Account, error) {
	args := m.Called()
	return args.Get(0).([]model.Account), args.Error(1)
}

func TestProcessAccountBalanceUpdates(t *testing.T) {
	mockConsumer := new(MockKafkaConsumer)
	mockProducer := new(MockKafkaProducer)
	mockAccountService := new(MockAccountService)

	processor := NewProcessor(mockConsumer, mockProducer, []string{"test-topic"}, []string{"status-topic"}, "test-group", mockAccountService)

	ctx := context.Background()

	mockConsumer.On("Consume", ctx, "test-topic", "test-group", mock.AnythingOfType("func(kafka.Message) error")).Return(nil).Run(func(args mock.Arguments) {
		handler := args.Get(3).(func(kafka.Message) error)
		message := kafka.Message{
			Key:   []byte("key"),
			Value: []byte(`{"accountId":"123", "amount":100.0, "transactionType":"credit"}`),
		}
		handler(message)
	})

	mockAccountService.On("UpdateAccountBalance", ctx, "123", 100.0, "credit").Return(nil)
	mockProducer.On("Produce", ctx, "status-topic", mock.AnythingOfType("kafka.Message")).Return(nil)

	err := processor.ProcessAccountBalanceUpdates(ctx)
	assert.NoError(t, err)

	mockConsumer.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
	mockAccountService.AssertExpectations(t)
}

func TestHandleAccountBalanceUpdate(t *testing.T) {
	mockProducer := new(MockKafkaProducer)
	mockAccountService := new(MockAccountService)

	processor := &processor{
		producer:       mockProducer,
		producerTopics: []string{"status-topic"},
		accountService: mockAccountService,
	}

	ctx := context.Background()
	message := kafka.Message{
		Key:   []byte("key"),
		Value: []byte(`{"accountId":"123", "amount":100.0, "transactionType":"credit"}`),
	}

	mockAccountService.On("UpdateAccountBalance", ctx, "123", 100.0, "credit").Return(nil)
	mockProducer.On("Produce", ctx, "status-topic", mock.AnythingOfType("kafka.Message")).Return(nil)

	err := processor.handleAccountBalanceUpdate(ctx, message)
	assert.NoError(t, err)

	mockAccountService.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestHandleAccountBalanceUpdate_Failure(t *testing.T) {
	mockProducer := new(MockKafkaProducer)
	mockAccountService := new(MockAccountService)

	processor := &processor{
		producer:       mockProducer,
		producerTopics: []string{"status-topic"},
		accountService: mockAccountService,
	}

	ctx := context.Background()
	message := kafka.Message{
		Key:   []byte("key"),
		Value: []byte(`{"accountId":"123", "amount":100.0, "transactionType":"credit"}`),
	}

	mockAccountService.On("UpdateAccountBalance", ctx, "123", 100.0, "credit").Return(errors.New("update error"))
	mockProducer.On("Produce", ctx, "status-topic", mock.AnythingOfType("kafka.Message")).Return(nil)

	err := processor.handleAccountBalanceUpdate(ctx, message)
	assert.NoError(t, err)

	mockAccountService.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}
