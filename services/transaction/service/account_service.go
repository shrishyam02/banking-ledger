package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/shrishyam02/banking-ledger/common/logger"
)

type accountService struct {
	AccountServiceURL string
}

type AccountService interface {
	GetAccountByID(ctx context.Context, accountID uuid.UUID) (map[string]any, error)
}

func NewAccountService(accountServiceURL string) AccountService {
	return &accountService{
		AccountServiceURL: accountServiceURL,
	}
}

func (s *accountService) GetAccountByID(ctx context.Context, accountID uuid.UUID) (map[string]any, error) {
	url := fmt.Sprintf("%s/api/v1/accounts/%s", s.AccountServiceURL, accountID.String())
	logger.Log.Info().Msgf("URL: %s", url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("test", "test")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("account not found")
	}

	var account map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, err
	}

	return account, nil
}
