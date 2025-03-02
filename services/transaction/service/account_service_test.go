package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAccountByID_Success(t *testing.T) {
	accountID := uuid.New()
	expectedAccount := map[string]any{"id": accountID.String(), "name": "Test Account"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/accounts/"+accountID.String(), r.URL.Path)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedAccount)
	}))
	defer server.Close()

	service := NewAccountService(server.URL)
	ctx := context.Background()

	account, err := service.GetAccountByID(ctx, accountID)
	assert.NoError(t, err)
	assert.Equal(t, expectedAccount, account)
}

func TestGetAccountByID_NotFound(t *testing.T) {
	accountID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	service := NewAccountService(server.URL)
	ctx := context.Background()

	account, err := service.GetAccountByID(ctx, accountID)
	assert.Error(t, err)
	assert.Nil(t, account)
	assert.Equal(t, "account not found", err.Error())
}

func TestGetAccountByID_RequestError(t *testing.T) {
	service := NewAccountService("http://invalid-url")
	ctx := context.Background()
	accountID := uuid.New()

	account, err := service.GetAccountByID(ctx, accountID)
	assert.Error(t, err)
	assert.Nil(t, account)
}

func TestGetAccountByID_DecodeError(t *testing.T) {
	accountID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	service := NewAccountService(server.URL)
	ctx := context.Background()

	account, err := service.GetAccountByID(ctx, accountID)
	assert.Error(t, err)
	assert.Nil(t, account)
}
