package repository

import (
	"account/model"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shrishyam02/banking-ledger/common/logger"
	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccount(account *model.Account) error
	GetAccountByID(id uuid.UUID) (*model.Account, error)
	ListAccounts() ([]model.Account, error)
	CreateOrUpdateCustomer(customer *model.Customer) error
	UpdateAccountBalance(ctx context.Context, accountID string, amount float64, transactionType string) error
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

func (r *accountRepository) CreateOrUpdateCustomer(customer *model.Customer) error {
	return r.db.Save(customer).Error
}

func (r *accountRepository) UpdateAccountBalance(ctx context.Context, accountID string, amount float64, transactionType string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var account model.Account
		if err := tx.First(&account, "ID = ?", accountID).Error; err != nil {
			return err
		}
		logger.Log.Info().Msgf("handleAccountBalanceUpdate account balance check update. %v account:(%v %v %v)", account, accountID, amount, transactionType)

		// Capture the current UpdatedAt timestamp
		currentUpdatedAt := account.UpdatedAt

		// Update the balance and set the new UpdatedAt timestamp
		if transactionType == "credit" {
			account.Balance = account.Balance + amount
		} else if transactionType == "debit" {
			account.Balance = account.Balance - amount
		} else {
			return fmt.Errorf("unknown transaction type %s", transactionType)
		}
		logger.Log.Info().Msgf("handleAccountBalanceUpdate account balance pre update. %v account:(%v %v %v)", account, accountID, amount, transactionType)

		//account.UpdatedAt = time.Now()

		// Use optimistic locking to ensure that the update is only applied if the UpdatedAt timestamp has not changed
		if err := tx.Model(&account).Where("UpdatedAt = ?", currentUpdatedAt).Updates(model.Account{Balance: account.Balance}).Error; err != nil {
			return err
		}

		// Check if the update was successful
		if tx.RowsAffected == 0 {
			return errors.New("concurrent update detected")
		}

		return nil
	})
}
