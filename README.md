## Transactions Service (Go)

This repository implements a simple banking transactions service following Hexagonal Architecture (Ports & Adapters / Clean Architecture). The goal is to demonstrate separation of concerns between domain, ports (interfaces), infrastructure adapters (Postgres), and delivery adapters (HTTP).

### Technologies
- Language: Go
- Database: PostgreSQL (init SQL in deploy/postgres/init.sql)
- Testing: go test (includes unit and integration tests)
- Docker / docker-compose for local environment

### Hexagonal Architecture Overview
- internal/core: domain (entities, errors) and ports (service/repository interfaces)
- internal/adapters: delivery adapters (HTTP handlers)
- internal/adapters/db/postgres: persistence adapters (repository implementations)
- cmd/app: application bootstrap (dependency injection and HTTP server)

### Main Structure
- cmd/app/main.go          -> entry point
- internal/core/...        -> domain, ports and services
- internal/adapters/...    -> HTTP handlers and Postgres repositories
- deploy/postgres/init.sql -> script for table creation and initial data
- internal/tests/...       -> unit and integration tests

### Configuration / Environment Variables
- DB_CONNECTION_STRING: Postgres connection string (required)
- HTTP_PORT: HTTP port (e.g., 8080)

### HTTP Endpoints
#### 1. Create Account
- POST /accounts
- Body: { "document_number": "<string>" }
- Responses:
  - 201 Created: { "id": int, "document_number": string, "created_at": string }
  - 400 Bad Request: invalid input
  - 409 Conflict: account already exists
  - Location header: /accounts/{id}

#### 2. Get Account by ID
- GET /accounts/{id}
- Responses:
  - 200 OK: { "id": int, "document_number": string, "created_at": string }
  - 400 Bad Request: invalid id
  - 404 Not Found: account not found

#### 3. Create Transaction
- POST /transactions
- Body: { "account_id": int, "operation_type_id": int, "amount": number }
- Responses:
  - 201 Created: { "id": int, "account_id": int, "operation_type_id": int, "amount": number, "event_date": string }
  - 400 Bad Request: invalid data (non-existent account, invalid operation, invalid amount)

### API Documentation (Swagger/OpenAPI)

The API has interactive documentation via Swagger UI:

- **Swagger UI**: http://localhost:5000/swagger/index.html
- **OpenAPI Specification (JSON)**: http://localhost:5000/swagger.json

Swagger documentation is automatically generated from comments in handler functions using `swaggo/swag`. Each endpoint includes:
- Description and summary
- Input parameters (path, query, body)
- Response schemas (success and error)
- HTTP status codes

#### How to Update Swagger Documentation
After making changes to handlers or adding new endpoints:

```bash
make swagger
# or
swag init -g cmd/app/main.go
```

### Request/Response Logging

All HTTP requests and responses are automatically logged in JSON. The logs include:

**Request Log Entry**:
```json
{
  "timestamp": "2026-04-13T22:37:01.258468357Z",
  "level": "INFO",
  "event": "http.request",
  "method": "POST",
  "path": "/accounts"
}
```

**Response Log Entry**:
```json
{
  "timestamp": "2026-04-13T22:37:01.258500000Z",
  "level": "INFO",
  "event": "http.response",
  "status_code": 201,
  "duration_ms": 15
}
```

#### Logging Features
- **Log Levels**: 
  - `INFO` for 2xx responses
  - `WARN` for 4xx responses  
  - `ERROR` for 5xx responses
- **Zero Performance Overhead**: Logging middleware adds <2ms latency per request

### Standardized Error Handling

All API errors follow a consistent JSON format via centralized error handling:

**Error Response Format**:
```json
{
  "error": "Not Found",
  "message": "Account not found",
  "code": "ACCOUNT_NOT_FOUND"
}
```

#### Error Handling Architecture

Error handling is centralized in `internal/adapters/handlers/error_handler.go`:

- **ErrorResponse**: Standardized response struct with error, message, and code
- **errorMapping**: Maps domain errors to HTTP status codes and error codes
- **RespondWithError()**: Converts domain errors to HTTP responses automatically
- **Helper functions**: `RespondWithBadRequest()`, `RespondWithNotFound()`, `RespondWithConflict()`, `RespondWithInternalServerError()`

This approach ensures:
- Consistent error format across all endpoints
- Centralized error code definitions (single source of truth)
- Easy to maintain and extend error handling
- Domain errors automatically mapped to correct HTTP status codes

### Notes on Operation Types
- The deploy/postgres/init.sql file inserts initial operation types (e.g., Normal Purchase, Withdrawal, Credit Voucher). Use the corresponding IDs when creating transactions.

### How to Run Locally
1. With Docker Compose (recommended):
   ```bash
   ./scripts/run.sh
   # or
   docker-compose up --build
   ```
   (the Postgres service will mount deploy/postgres/init.sql)

2. Run locally without Docker (Postgres already available):
   ```bash
   export DB_CONNECTION_STRING="postgresql://user:pass@host:5432/dbname?sslmode=disable"
   export HTTP_PORT=8080
   go run ./cmd/app
   ```

### Examples with curl
- Create account:
```bash
  curl -v -X POST http://localhost:8080/accounts \
    -H "Content-Type: application/json" \
    -d '{"document_number":"12345678900"}'
```

- Create transaction:
```bash
  curl -v -X POST http://localhost:8080/transactions \
    -H "Content-Type: application/json" \
    -d '{"account_id":1,"operation_type_id":1,"amount":100.5}'
```

### Testing
- Run full test suite:
  ```bash
  go test ./...
  ```

- Run specific integration tests (e.g., using Testcontainers/Docker):
  ```bash
  go test ./internal/tests/integration/handlers/http -run TestCreateAccount -v
  ```
