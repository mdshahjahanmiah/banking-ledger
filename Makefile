.PHONY: build run ledger processor test unit-test feature-test

build:
	go build -o bin/transaction_ledger ./cmd/transaction_ledger
	go build -o bin/transaction_processor ./cmd/transaction_processor

run-ledger: build
	./bin/transaction_ledger

run-processor: build
	./bin/transaction_processor

run-all: run-ledger run-processor

unit-test:
	go test -v ./... -count=1

feature-test:
	go test -v -count=1 ./... -tags=feature
lint:
	golangci-lint run ./...

test: unit-test feature-test
