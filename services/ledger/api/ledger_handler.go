package api

import (
	"ledger/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ledgerHandler struct {
	service service.LedgerService
}

type LedgerHandler interface {
	GetAccountTransactionHistory(c *gin.Context)
	GetTransactionHistory(c *gin.Context)
}

func NewledgerHandler(service service.LedgerService) LedgerHandler {
	return &ledgerHandler{service: service}
}

func (h *ledgerHandler) GetAccountTransactionHistory(c *gin.Context) {
	accountID := c.Param("id")
	transactions, err := h.service.GetAccountTransactionHistory(c.Request.Context(), accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

func (h *ledgerHandler) GetTransactionHistory(c *gin.Context) {
	accountID := c.Param("id")
	transactions, err := h.service.GetTransactionHistory(c.Request.Context(), accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}
