package api

import (
	"net/http"

	"account/model"
	"account/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type accountHandler struct {
	service service.AccountService
}

type AccountHandler interface {
	CreateAccount(c *gin.Context)
	GetAccount(c *gin.Context)
	ListAccounts(c *gin.Context)
}

func NewAccountHandler(service service.AccountService) AccountHandler {
	return &accountHandler{service: service}
}

func (h *accountHandler) CreateAccount(c *gin.Context) {
	var account model.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if account.CustomerID == uuid.Nil {
		account.CustomerID = uuid.New()
	}
	if err := h.service.CreateOrUpdateCustomer(&account.Customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	account.Status = "active"
	if err := h.service.CreateAccount(&account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}

func (h *accountHandler) GetAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}
	account, err := h.service.GetAccountByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	c.JSON(http.StatusOK, account)
}

func (h *accountHandler) ListAccounts(c *gin.Context) {
	accounts, err := h.service.ListAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}
