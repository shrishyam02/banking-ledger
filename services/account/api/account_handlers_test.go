package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"account/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAccountService struct {
	mock.Mock
}

func (m *MockAccountService) CreateAccount(account *model.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountService) GetAccountByID(id uuid.UUID) (*model.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockAccountService) ListAccounts() ([]model.Account, error) {
	args := m.Called()
	return args.Get(0).([]model.Account), args.Error(1)
}

func (m *MockAccountService) CreateOrUpdateCustomer(customer *model.Customer) error {
	args := m.Called(customer)
	return args.Error(0)
}

func (m *MockAccountService) UpdateAccountBalance(ctx context.Context, id string, balance float64, currency string) error {
	args := m.Called(ctx, id, balance, currency)
	return args.Error(0)
}

func TestCreateAccount(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/accounts", handler.CreateAccount)

	account := model.Account{ID: uuid.New(), AccountNumber: "12345"}
	mockService.On("CreateAccount", &account).Return(nil)

	body, _ := json.Marshal(account)
	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	mockService.AssertExpectations(t)
}

func TestGetAccount(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/accounts/:id", handler.GetAccount)

	account := model.Account{ID: uuid.New(), AccountNumber: "123456"}
	mockService.On("GetAccountByID", account.ID).Return(&account, nil)

	req, _ := http.NewRequest(http.MethodGet, "/accounts/"+account.ID.String(), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestListAccounts(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/accounts", handler.ListAccounts)

	accounts := []model.Account{
		{ID: uuid.New(), AccountNumber: "123456"},
		{ID: uuid.New(), AccountNumber: "123456"},
	}
	mockService.On("ListAccounts").Return(accounts, nil)

	req, _ := http.NewRequest(http.MethodGet, "/accounts", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}
