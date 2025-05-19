# Banking Ledger

A banking ledger service built for transaction processing and ledger management.

##  Features

- Support the creation of accounts with specified initial balances.
- Facilitate deposits and withdrawals
- Maintain a detailed transaction log (ledger) for each account
- PostgreSQL for transaction data storage, and MongoDB for additional data persistence
- Ensured ACID-like consistency for core operations to prevent double spending or inconsistent balances
- Integrated an asynchronous queue or broker to manage transaction requests efficiently
- Consumer connect, ping and automatic reconnect when available
- Unit and feature test coverage with BDD (Behavior Driven Development)
- Docker containerization for easy deployment

## Architecture

The service consists of two main components:

1. **Transaction Ledger Service**
   - REST API server (port 3000)
   - Handles transaction requests
   - Manages ledger entries
   - Communicates with Kafka for event publishing

2. **Transaction Processor**
   - Processes transaction events from Kafka
   - Updates ledger entries
   - Handles transaction state management

## Prerequisites

- Docker and Docker Compose
- Go 1.23.5 or later

### Kafka Topic Initialization

During initialization, the following two Kafka topics will be created:

1. `transactions` - Main topic where all valid transaction events are published.
2. `transactions-dlq` - Dead Letter Queue for storing failed or unprocessable transaction messages.

---

### Verify Topic Creation

To verify that the topics were successfully created, run:

```bash
docker exec -it kafka kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --list
  ```
### If Topics Are Missing
In some cases, the topics may not be created automatically. You can manually verify or create them using the following commands:

Describe topic to check if it exists:
```bash
docker exec -it kafka kafka-topics.sh \
--bootstrap-server kafka:9092 \
--describe --topic transactions-dlq
  ```

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/mdshahjahanmiah/banking-ledger.git
   cd banking-ledger
   ```

2. Start the services using Docker Compose:
   ```bash
   docker-compose up -d
   ```

3. Run database migrations:
   ```bash
   make migrate-up
   ```

4. Run tests:
   ```bash
   make test
   ```

## Service Configuration

The services are configured through environment variables:

- `POSTGRES_DSN`: PostgreSQL connection string
- `KAFKA_BROKER_URL`: Kafka broker address
- `MONGO_URI`: MongoDB connection string

Default values are set in the `docker-compose.yml` file.

## Testing

The project uses multiple testing approaches:

- BDD tests using Godog
- Unit tests with Go's testing package
- Integration tests with sqlmock
- API tests

Run tests using:
```bash
make test        # Run all tests
make test-bdd    # Run BDD tests
make test-unit   # Run unit tests
```

## Future Improvements & Known Issues

### Planned Improvements
- Implement connection pooling for database connections
- Add caching layer using Redis for frequently accessed data
- Optimize Kafka consumer group configurations
- Add rate limiting for API endpoints
- Implement API key management
- Add request validation middleware
- Implement distributed tracing with OpenTelemetry
- Enhance logging with structured logging
- Add Grafana dashboards for monitoring
- Add support for batch transactions
- Implement transaction reconciliation process

### Known Issues
- Some test cases need better coverage
- Documentation needs more examples
- Some hardcoded values need to be moved to configuration
