// Package consumer provides Kafka consumer functionality for processing financial transactions.
//
// It handles the consumption of messages from the "transactions" topic,
// delegates transaction processing to the service layer, applies retry logic,
// and forwards failed messages to a Dead Letter Queue (DLQ). It also
// persists all processed transactions to an audit repository for traceability.
//
// This package is intended to be resilient and observant, ensuring durable
// and recoverable transaction ingestion in a distributed system.
package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/mdshahjahanmiah/banking-ledger/model"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/config"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/transaction"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/mdshahjahanmiah/explore-go/repository"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	Logger             *logging.Logger
	Reader             *kafka.Reader
	DLTWriter          *kafka.Writer
	TransactionService transaction.Service
	AuditRepo          *repository.Repository[model.Transaction]
	Config             config.Config
}

func NewConsumer(cfg config.Config, logger *logging.Logger, service transaction.Service, auditRepo *repository.Repository[model.Transaction]) *Consumer {
	return &Consumer{
		Logger:             logger,
		Config:             cfg,
		TransactionService: service,
		AuditRepo:          auditRepo,
		DLTWriter: &kafka.Writer{
			Addr:     kafka.TCP(cfg.KafkaBrokerURL),
			Topic:    "transactions-dlq",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

// Start begins consuming messages from Kafka and processes them.
// It handles retries, permanent failures, audit logging, and dead-letter routing.
func (c *Consumer) Start(errorChan chan error, doneChan chan struct{}) {
	go func() {
		defer func() {
			if c.Reader != nil {
				_ = c.Reader.Close()
			}
			close(doneChan)
		}()

		// Retry until initial connection is established
		for {
			if err := c.connect(); err != nil {
				c.Logger.Error("Initial Kafka connection failed", "error", err)
				time.Sleep(5 * time.Second)
				continue
			}
			break
		}

		c.Logger.Info("Kafka connected, starting message consumption...")

		for {
			// Read message from Kafka
			msg, err := c.Reader.ReadMessage(context.Background())
			if err != nil {
				if errors.Is(err, context.Canceled) {
					c.Logger.Info("Kafka consumer context canceled, shutting down...")
					return
				}

				c.Logger.Error("Kafka read error", "error", err)
				c.reconnect()
				continue
			}

			c.handleMessage(msg)
		}
	}()
}

// handleMessage processes a single Kafka message. It unmarshal the payload into a Transaction model,
// attempts to process the transaction with retries, writes to DLQ on failure, and audits the result.
func (c *Consumer) handleMessage(msg kafka.Message) {
	var txn model.Transaction
	if err := json.Unmarshal(msg.Value, &txn); err != nil {
		c.Logger.Error("Invalid transaction format", "error", err)
		return
	}

	const maxRetries = 3
	var attempt int
	var lastErr error

	for attempt = 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		lastErr = c.TransactionService.ProcessTransaction(ctx, txn)
		cancel()

		if lastErr == nil {
			txn.Status = transaction.TransactionStatusCompleted
			break
		}

		if errors.Is(lastErr, transaction.ErrDuplicateTransaction) ||
			errors.Is(lastErr, transaction.ErrAccountNotFound) {
			c.Logger.Warn("Permanent transaction failure, skipping retry", "id", txn.ID, "error", lastErr)
			txn.Status = transaction.TransactionStatusFailed
			lastErr = nil // clear error to avoid DLQ
			break
		}

		c.Logger.Warn("Retryable transaction failure", "id", txn.ID, "attempt", attempt, "error", lastErr)
		time.Sleep(time.Second * time.Duration(attempt)) // simple backoff
	}

	if lastErr != nil && txn.Status != transaction.TransactionStatusCompleted {
		txn.Status = transaction.TransactionStatusFailed
		payload, _ := json.Marshal(map[string]interface{}{
			"transaction": txn,
			"error":       lastErr.Error(),
			"failedAt":    time.Now(),
		})

		err := c.DLTWriter.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(txn.ID),
			Value: payload,
		})
		if err != nil {
			c.Logger.Error("Failed to write to DLT", "id", txn.ID, "error", err)
		} else {
			c.Logger.Info("Message sent to DLT", "id", txn.ID)
		}
	}

	c.Logger.Info("Transaction processed", "id", txn.ID, "status", txn.Status, "amount", txn.Amount)
	if err := c.AuditRepo.Save(txn); err != nil {
		c.Logger.Error("Audit failed (non-critical)", "id", txn.ID, "error", err)
	}

	c.Logger.Info("Transaction processed", "id", txn.ID, "status", txn.Status, "duration", time.Since(txn.CreatedAt))
}

// ping verifies Kafka broker availability by checking partition metadata.
func (c *Consumer) ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := kafka.DialContext(ctx, "tcp", c.Config.KafkaBrokerURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
		return err
	}

	_, err = conn.ReadPartitions()
	return err
}

// connect initializes the Kafka reader and logs discovered partitions.
func (c *Consumer) connect() error {
	if c.Reader != nil {
		_ = c.Reader.Close()
	}

	c.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{c.Config.KafkaBrokerURL},
		Topic:   "transactions",
		GroupID: "transaction-processor",
	})

	conn, err := kafka.Dial("tcp", c.Config.KafkaBrokerURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return err
	}
	c.Logger.Info("Kafka partitions found", "count", len(partitions))

	return nil
}

// reconnect attempts to re-establish Kafka connection with exponential backoff.
func (c *Consumer) reconnect() {
	backoff := time.Second

	for {
		c.Logger.Warn("Pinging Kafka for availability...")
		if err := c.ping(); err != nil {
			c.Logger.Error("Kafka unavailable", "error", err)
			time.Sleep(backoff)

			if backoff < 30*time.Second {
				backoff *= 2
			}
			continue
		}

		c.Logger.Info("Kafka is online, reconnecting reader...")
		if err := c.connect(); err != nil {
			c.Logger.Error("Reconnect failed", "error", err)
			time.Sleep(backoff)
			continue
		}

		c.Logger.Info("Kafka successfully reconnected.")
		return
	}
}
