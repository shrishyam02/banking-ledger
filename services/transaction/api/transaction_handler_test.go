package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"transaction/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKafkaWriter struct {
	mock.Mock
}

func (m *MockKafkaWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	args := m.Called(ctx, msgs)
	return args.Error(0)
}

func (m *MockKafkaWriter) Produce(ctx context.Context, topic string, msg kafka.Message) error {
	args := m.Called(ctx, topic, msg)
	return args.Error(0)
}

type MockAccountService struct {
	mock.Mock
}

func (m *MockAccountService) GetAccountByID(ctx context.Context, accountID uuid.UUID) (map[string]interface{}, error) {
	args := m.Called(ctx, accountID)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func TestCreateTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockKafkaWriter := new(MockKafkaWriter)
	mockAccountService := new(MockAccountService)
	handler := NewTransactionHandler(mockKafkaWriter, []string{"topic1"}, mockAccountService)

	router := gin.Default()
	router.POST("/transactions", handler.CreateTransaction)

	t.Run("should return 400 if request body is invalid", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer([]byte(`invalid`)))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should return 404 if account is not found", func(t *testing.T) {
		transaction := model.Transaction{AccountID: uuid.New()}
		body, _ := json.Marshal(transaction)
		req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()

		mockAccountService.On("GetAccountByID", mock.Anything, transaction.AccountID).Return(nil, errors.New("account not found"))

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("should return 400 if account is not active", func(t *testing.T) {
		transaction := model.Transaction{AccountID: uuid.New()}
		body, _ := json.Marshal(transaction)
		req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()

		mockAccountService.On("GetAccountByID", mock.Anything, transaction.AccountID).Return(map[string]interface{}{"Status": "inactive"}, nil)

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should return 500 if kafka writer fails", func(t *testing.T) {
		transaction := model.Transaction{AccountID: uuid.New()}
		body, _ := json.Marshal(transaction)
		req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()

		mockAccountService.On("GetAccountByID", mock.Anything, transaction.AccountID).Return(map[string]interface{}{"Status": "active"}, nil)
		mockKafkaWriter.On("WriteMessages", mock.Anything, mock.Anything).Return(errors.New("kafka error"))

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("should return 201 if transaction is created successfully", func(t *testing.T) {
		transaction := model.Transaction{AccountID: uuid.New()}
		body, _ := json.Marshal(transaction)
		req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()

		mockAccountService.On("GetAccountByID", mock.Anything, transaction.AccountID).Return(map[string]interface{}{"Status": "active"}, nil)
		mockKafkaWriter.On("WriteMessages", mock.Anything, mock.Anything).Return(nil)

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
	})
}
