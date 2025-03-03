package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLedgerService struct {
	mock.Mock
}

func (m *MockLedgerService) HandleMessage(ctx context.Context, msg kafka.Message) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *MockLedgerService) GetAccountTransactionHistory(ctx context.Context, accountID string) ([]map[string]interface{}, error) {
	args := m.Called(ctx, accountID)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockLedgerService) GetTransactionHistory(ctx context.Context, accountID string) ([]map[string]interface{}, error) {
	args := m.Called(ctx, accountID)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func TestGetAccountTransactionHistory(t *testing.T) {
	mockService := new(MockLedgerService)
	handler := NewledgerHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/account/:id/transactions", handler.GetAccountTransactionHistory)

	t.Run("success", func(t *testing.T) {
		mockTransactions := []map[string]interface{}{{"id": "1"}, {"id": "2"}}
		mockService.On("GetAccountTransactionHistory", mock.Anything, "123").Return(mockTransactions, nil)

		req, _ := http.NewRequest(http.MethodGet, "/account/123/transactions", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockService.On("GetAccountTransactionHistory", mock.Anything, "123").Return(nil, errors.New("some error"))

		req, _ := http.NewRequest(http.MethodGet, "/account/123/transactions", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetTransactionHistory(t *testing.T) {
	mockService := new(MockLedgerService)
	handler := NewledgerHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/transactions/:id", handler.GetTransactionHistory)

	t.Run("success", func(t *testing.T) {
		mockTransactions := []map[string]interface{}{{"id": "1"}, {"id": "2"}}
		mockService.On("GetTransactionHistory", mock.Anything, "123").Return(mockTransactions, nil)

		req, _ := http.NewRequest(http.MethodGet, "/transactions/123", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockService.On("GetTransactionHistory", mock.Anything, "123").Return(nil, errors.New("some error"))

		req, _ := http.NewRequest(http.MethodGet, "/transactions/123", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockService.AssertExpectations(t)
	})
}
