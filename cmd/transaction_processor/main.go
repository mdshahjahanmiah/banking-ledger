package main

import (
	"github.com/mdshahjahanmiah/banking-ledger/cmd/transaction_processor/consumer"
	"github.com/mdshahjahanmiah/banking-ledger/model"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/broker"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/config"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/db"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/transaction"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/mdshahjahanmiah/explore-go/repository"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
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

	producer := broker.NewKafkaProducer(cfg.KafkaBrokerURL, "")
	if err != nil {
		logger.Error("failed to create Kafka producer", "err", err)
		return
	}

	// Prepare transaction store
	auditRepo := repository.NewMongoRepository[model.Transaction](mongoDB.Client, "ledger", "transactions")
	txnService, err := transaction.NewService(cfg, logger, database, auditRepo, producer)
	if err != nil {
		logger.Error("failed to initialize service", "err", err)
		return
	}

	processor := consumer.NewConsumer(cfg, logger, txnService, auditRepo)

	errorChan := make(chan error)
	doneChan := make(chan struct{})

	// Start the processor
	logger.Info("consumer.Start() called")
	processor.Start(errorChan, doneChan)

	// Optional: handle shutdown signals (CTRL+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errorChan:
		logger.Error("processor error", "error", err)
	case sig := <-sigChan:
		logger.Info("received shutdown signal", "signal", sig)
	case <-doneChan:
		logger.Info("processor completed work and exited")
	}

	// Cleanup
	logger.Info("shutting down gracefully...")
}
