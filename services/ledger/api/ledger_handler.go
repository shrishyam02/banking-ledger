package api

import (
	"ledger/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LedgerHandler struct {
	service *service.LedgerService
}

func NewLedgerHandler(service *service.LedgerService) *LedgerHandler {
	return &LedgerHandler{service: service}
}

func (h *LedgerHandler) GetTransactionHistory(c *gin.Context) {
	accountID := c.Param("id")
	transactions, err := h.service.GetTransactionHistory(c.Request.Context(), accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}
