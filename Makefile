
start: docs
	go run cmd/app/main.go

.PHONY: docs install-mockgen mock
docs:
	swag init -g cmd/app/main.go -o ./docs/swagger/

install-mockgen:
	go install go.uber.org/mockgen@latest

mock:
	mockgen -source=internal/config/config.go -destination=./files/mocks/config/mock_config.go
	mockgen -source=internal/module/health/health.go -destination=./files/mocks/health/mock_health.go
	mockgen -source=internal/module/workflow/workflow.go -destination=./files/mocks/workflow/mock_workflow.go

