.PHONY: build run ledger processor

build:
	go build -o bin/transaction_ledger ./cmd/transaction_ledger
	go build -o bin/transaction_processor ./cmd/transaction_processor

run-ledger: build
	./bin/transaction_ledger

run-processor: build
	./bin/transaction_processor

run-all: run-ledger run-processor
