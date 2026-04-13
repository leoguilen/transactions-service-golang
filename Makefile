.PHONY: build test run

build:
	go build -o bin/transactions ./cmd/app

test:
	go test ./internal/tests/../... -v

run:
	go run ./cmd/app

run-dev:
	DB_CONNECTION_STRING="postgres://postgres:postgres@localhost:5432/transactions_db?sslmode=disable" HTTP_PORT=5000 go run ./cmd/app

run-docker:
	chmod +x ./scripts/run.sh
	./scripts/run.sh