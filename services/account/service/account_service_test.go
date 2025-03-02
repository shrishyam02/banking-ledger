package service

import (
	"account/model"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAccountRepository is a mock implementation of the AccountRepository interface
type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) CreateAccount(account *model.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountRepository) GetAccountByID(id uuid.UUID) (*model.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockAccountRepository) ListAccounts() ([]model.Account, error) {
	args := m.Called()
	return args.Get(0).([]model.Account), args.Error(1)
}

func TestCreateAccount(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	service := NewService(mockRepo)

	account := &model.Account{ID: uuid.New()}
	mockRepo.On("CreateAccount", account).Return(nil)

	err := service.CreateAccount(account)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAccountByID(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	service := NewService(mockRepo)

	accountID := uuid.New()
	account := &model.Account{ID: accountID}
	mockRepo.On("GetAccountByID", accountID).Return(account, nil)

	result, err := service.GetAccountByID(accountID)
	assert.NoError(t, err)
	assert.Equal(t, account, result)
	mockRepo.AssertExpectations(t)
}

func TestListAccounts(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	service := NewService(mockRepo)

	accounts := []model.Account{
		{ID: uuid.New()},
		{ID: uuid.New()},
	}
	mockRepo.On("ListAccounts").Return(accounts, nil)

	result, err := service.ListAccounts()
	assert.NoError(t, err)
	assert.Equal(t, accounts, result)
	mockRepo.AssertExpectations(t)
}

func (m *MockAccountRepository) CreateOrUpdateCustomer(customer *model.Customer) error {
	args := m.Called(customer)
	return args.Error(0)
}

func (m *MockAccountRepository) UpdateAccountBalance(ctx context.Context, id string, balance float64, currency string) error {
	args := m.Called(ctx, id, balance, currency)
	return args.Error(0)
}

func TestCreateAccount_Error(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	service := NewService(mockRepo)

	account := &model.Account{ID: uuid.New()}
	mockRepo.On("CreateAccount", account).Return(errors.New("create error"))

	err := service.CreateAccount(account)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
