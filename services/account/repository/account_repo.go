package repository

import (
	"account/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccount(account *model.Account) error
	GetAccountByID(id uuid.UUID) (*model.Account, error)
	ListAccounts() ([]model.Account, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) CreateAccount(account *model.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepository) GetAccountByID(id uuid.UUID) (*model.Account, error) {
	var account model.Account
	err := r.db.Preload("Customer").First(&account, id).Error
	return &account, err
}

func (r *accountRepository) ListAccounts() ([]model.Account, error) {
	var accounts []model.Account
	err := r.db.Preload("Customer").Find(&accounts).Error
	return accounts, err
}
