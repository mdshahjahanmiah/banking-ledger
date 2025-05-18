FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY . .
RUN go build -o transaction_ledger ./cmd/transaction_ledger/main.go && \
    go build -o transaction_processor ./cmd/transaction_processor/main.go

# Expose HTTP port (used by transaction_ledger)
EXPOSE 3000
