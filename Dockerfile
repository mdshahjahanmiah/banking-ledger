FROM golang:1.23

WORKDIR /app

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy rest and build both binaries
COPY . .
RUN go build -o transaction_ledger ./cmd/transaction_ledger/main.go
RUN go build -o transaction_processor ./cmd/transaction_processor/main.go

EXPOSE 3000
