APP_NAME=engine
USER=root
PASSWORD=root
DB_NAME=sequence_db
DB_STRING=postgres://"$(USER)":"$(PASSWORD)"@localhost:5432/"$(DB_NAME)"?sslmode=disable
KAFKA_CONTAINER=kafka
KAFKA_BIN=/usr/bin/kafka-topics

.PHONY: run build test clean migrate swagger fmt

# Run the app
run:
	go run cmd/app/main.go

# Build binary
build:
	go build -o $(APP_NAME) cmd/app/main.go

# Docker
docker-build:
	docker build -t sequence-service:latest .

up:
	docker compose up -d

down:
	docker compose down

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

# Generate Mocks
mock:
	mockgen -source=internal/config/config.go -destination=./files/mocks/config/mock_config.go
	mockgen -source=internal/module/health/health.go -destination=./files/mocks/health/mock_health.go
	mockgen -source=internal/module/workflow/workflow.go -destination=./files/mocks/workflow/mock_workflow.go

# Create Kafka topics
kafka-topics:
	docker exec -it $(KAFKA_CONTAINER) $(KAFKA_BIN) --bootstrap-server localhost:9092 --create --topic email-jobs --partitions 1 --replication-factor 1
	docker exec -it $(KAFKA_CONTAINER) $(KAFKA_BIN) --bootstrap-server localhost:9092 --create --topic followup-events --partitions 1 --replication-factor 1
	docker exec -it $(KAFKA_CONTAINER) $(KAFKA_BIN) --bootstrap-server localhost:9092 --create --topic email-retries --partitions 1 --replication-factor 1
	docker exec -it $(KAFKA_CONTAINER) $(KAFKA_BIN) --bootstrap-server localhost:9092 --create --topic email-events --partitions 1 --replication-factor 1

# Consume messages from a topic
kafka-consume:
	@read -p "Enter topic name: " topic; \
	docker exec -it kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic email-jobs --from-beginning
	
# Push message to Kafka topic
kafka-produce:
	@read -p "Enter topic name: " TOPIC; \
	read -p "Enter message: " MSG; \
	docker exec -i kafka kafka-console-producer \
		--bootstrap-server localhost:9092 \
		--topic $$TOPIC <<< "$$MSG"

