FROM golang:1.23

WORKDIR /app

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy binaries
COPY . .
RUN go build -o transaction_ledger ./cmd/transaction_ledger/main.go

EXPOSE 3000
