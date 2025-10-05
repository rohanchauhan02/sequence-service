APP_NAME=engine
USER=root
PASSWORD=root
DB_NAME=sequence_db
DB_STRING=postgres://"$(USER)":"$(PASSWORD)"@localhost:5432/"$(DB_NAME)"?sslmode=disable

.PHONY: run build test clean migrate swagger fmt

# Run the app
run:
	go run cmd/app/main.go

# Build binary
build:
	go build -o $(APP_NAME) cmd/app/main.go

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
	swag init -g cmd/app/main.go -o docs/swagger/

# Code formatting
fmt:
	go fmt ./...

