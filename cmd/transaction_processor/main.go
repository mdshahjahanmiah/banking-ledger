package main

import (
	"context"
	"encoding/json"
	"github.com/mdshahjahanmiah/banking-ledger/repository"
	"log/slog"
	"time"

	"github.com/mdshahjahanmiah/banking-ledger/model"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/config"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/db"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/transaction"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/segmentio/kafka-go"
)

func main() {
	slog.Info("transaction processor is starting...")

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		return
	}
	slog.Info("config loaded", "kafka", cfg.KafkaBrokerURL)

	logger, err := logging.NewLogger(cfg.LoggerConfig)
	if err != nil {
		slog.Error("failed to initialize logger", "err", err)
		return
	}

	// Connect to SQL DB
	database, err := db.NewDB(cfg.PostgresDSN, logger)
	if err != nil {
		logger.Error("failed to connect to database", "err", err)
		return
	}
	defer database.Close()

	// Connect to MongoDB
	mongoDB, err := db.NewMongoDB(cfg)
	if err != nil {
		logger.Error("MongoDB connection failed", "err", err)
		return
	}
	defer mongoDB.Close()

	// Prepare transaction store
	txnStore := transaction.NewStore(database)
	auditRepo := repository.NewMongoRepository[model.Transaction](mongoDB.Client, "ledger", "transactions")

	// Kafka consumer
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.KafkaBrokerURL},
		Topic:   "transactions",
		GroupID: "transaction-processor",
	})

	logger.Info("transaction processor started, waiting for Kafka messages...")
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			logger.Error("error reading message", "error", err)
			continue
		}

		var txn model.Transaction
		if err := json.Unmarshal(msg.Value, &txn); err != nil {
			logger.Error("invalid transaction format", "error", err)
			continue
		}

		// Create processing context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Process transaction
		if err := txnStore.ProcessTransaction(ctx, txn); err != nil {
			logger.Error("transaction processing failed", "id", txn.ID, "status", txn.Status, "error", err)
			txn.Status = transaction.TransactionStatusFailed
		} else {
			txn.Status = transaction.TransactionStatusCompleted
		}

		// Audit regardless of status
		if err := auditRepo.Save(txn); err != nil {
			logger.Error("audit failed", "id", txn.ID, "error", err)
		}

		logger.Info("transaction processed", "id", txn.ID, "status", txn.Status, "duration", time.Since(txn.CreatedAt))
	}
}
