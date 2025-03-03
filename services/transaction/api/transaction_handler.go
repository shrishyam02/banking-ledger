package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"transaction/model"
	"transaction/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	ckafka "github.com/shrishyam02/banking-ledger/common/kafka"
	"github.com/shrishyam02/banking-ledger/common/logger"
)

type transactionHandler struct {
	kafkaProducer  ckafka.KafkaProducer
	producerTopics []string
	accountService service.AccountService
}

type TransactionHandler interface {
	CreateTransaction(c *gin.Context)
}

func NewTransactionHandler(kafkaProducer ckafka.KafkaProducer, producerTopics []string, accountService service.AccountService) TransactionHandler {
	return &transactionHandler{
		kafkaProducer:  kafkaProducer,
		producerTopics: producerTopics,
		accountService: accountService,
	}
}

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var transaction model.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.accountService.GetAccountByID(c, transaction.AccountID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if status, ok := account["Status"].(string); !ok || status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account is not active"})
		return
	}

	transaction.ID = uuid.New()
	transaction.AcceptedAt = time.Now().UTC()

	transactionBytes, err := json.Marshal(transaction)
	if err != nil {
		logger.Log.Error().Msgf("Failed to marshal json message. err: %v", err)
		return
	}

	message := kafka.Message{
		Key:   []byte(transaction.ID.String()),
		Value: transactionBytes,
	}
	ctx := context.Background()
	kerr := h.kafkaProducer.Produce(ctx, h.producerTopics[0], message)

	if kerr != nil {
		logger.Log.Error().Msgf("Failed to publish message to kafka. err: %v", kerr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process transaction"})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}
