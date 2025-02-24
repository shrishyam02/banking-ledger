package service

import (
	"account/model"
	"account/repository"

	"github.com/google/uuid"
)

type AccountService interface {
	CreateAccount(account *model.Account) error
	GetAccountByID(id uuid.UUID) (*model.Account, error)
	ListAccounts() ([]model.Account, error)
}

type accountService struct {
	repo repository.AccountRepository
}

func NewService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (s *accountService) CreateAccount(account *model.Account) error {
	return s.repo.CreateAccount(account)
}

func (s *accountService) GetAccountByID(id uuid.UUID) (*model.Account, error) {
	return s.repo.GetAccountByID(id)
}

func (s *accountService) ListAccounts() ([]model.Account, error) {
	return s.repo.ListAccounts()
}
