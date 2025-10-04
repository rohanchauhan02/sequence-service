APP_NAME=sequence-service
DB_STRING=postgres://user:pass@localhost:5432/sequence_db?sslmode=disable

.PHONY: run build test clean migrate swagger fmt lint

# Run the app
run:
	go run cmd/api/main.go

# Build binary
build:
	go build -o bin/$(APP_NAME) cmd/api/main.go

# Run all tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out

# Database migrations
migrate-up:
	goose -dir database/migrations postgres "$(DB_STRING)" up

migrate-down:
	goose -dir database/migrations postgres "$(DB_STRING)" down

migrate-status:
	goose -dir database/migrations postgres "$(DB_STRING)" status

migrate-create:
	goose -dir database/migrations postgres "$(DB_STRING)" create $(name) sql

# Generate Swagger docs
swagger:
	swag init -g cmd/api/main.go -o docs

# Code formatting
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run
