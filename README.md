# Sequence Service

A high-performance email sequencing platform built with Go and PostgreSQL/Kafka. Manage multi-step email campaigns with intelligent scheduling.

---

## 🚀 Features

* **Sequence Management**: Create multi-step email sequences
* **Contact Management**: Add contacts to sequences
* **Email Tracking**: Open and click tracking
* **RESTful API**: Clean, documented endpoints with Swagger
* **PostgreSQL + Kafka**: Reliable data storage & messaging
* **Makefile + Docker Compose**: Easy local development

---

## 🛠 Tech Stack

* **Go 1.23+** - Backend language
* **Echo** - HTTP framework
* **PostgreSQL** - Database
* **Kafka** - Messaging
* **GORM** - ORM
* **Viper** - Configuration
* **Goose** - Migrations
* **Swagger** - API documentation

---

## ⚡ Quick Start (Docker Compose)

### 1. Clone & Setup

```bash
git clone https://github.com/rohanchauhan02/sequence-service
cd sequence-service
cp configs/app.config.sample.yml configs/app.config.local.yml
```

### 2. Update Configuration

Edit `configs/app.config.local.yml` (if needed):

```yaml
PORT: 8080

DB:
  HOST: localhost
  PORT: 5432
  NAME: sequence_db
  USER: your_username
  PASSWORD: your_password
  SSL_MODE: disable

KAFKA:
  BROKERS: kafka:9092
  TOPICS:
    EMAIL_JOBS: email-jobs
    FOLLOWUP_EVENTS: followup-events
    EMAIL_RETRIES: email-retries
    EMAIL_EVENTS: email-events
```

> Note: Hostnames `postgres` and `kafka` match the Docker Compose service names.

### 3. Run the Entire Stack

```bash
make up
```

This command will:

1. Start **PostgreSQL**
2. Start **Kafka + Zookeeper**
3. Build and run the **Sequence Service**

Server: **[http://localhost:8080](http://localhost:8080)**
Swagger UI: **[http://localhost:8080/swagger](http://localhost:8080/swagger)**

### 4. Stop the Stack

```bash
make down
```

---

## 🛠 Development Commands (Makefile)

```bash
# Run locally (Go only, requires local DB/Kafka)
make run

# Build binary
make build

# Run tests
make test

# Database migrations
make migrate-up
make migrate-down
make migrate-status

# Generate Swagger docs
make swagger

# Format code
make fmt

# Clean build artifacts
make clean

# Docker
make docker-build
make up
make down
```

---

## 📚 API Documentation

Visit **[http://localhost:8080/swagger](http://localhost:8080/swagger)** after starting the service.

### Example Endpoints

#### Health Check

```http
GET /api/v1/health
```

#### Create Sequence

```http
POST /api/v1/sequences
Content-Type: application/json

{
  "name": "Welcome Sequence",
  "open_tracking_enabled": true,
  "click_tracking_enabled": true,
  "steps": [
    {
      "step_order": 0,
      "subject": "Welcome!",
      "content": "Hi {name}, welcome!",
      "wait_days": 0
    }
  ]
}
```

---

## 🗄 Database

### Main Tables

* `sequences` - Sequence definitions
* `steps` - Email steps
* `mailboxes` - Email accounts
* `contacts` - Recipients
* `sequence_contacts` - Links contacts to sequences
* `email_queues` - Scheduled emails

### Migration Commands

```bash
make migrate-up          # Run migrations
make migrate-down        # Rollback last migration
make migrate-status      # Check migration status
```

---

## 🏗 Project Structure

```
sequence-service/
├── Makefile                    # Development commands
├── cmd/api/main.go            # App entry point
├── docs/                      # Swagger documentation
├── internal/
│   ├── app/app.go             # App setup
│   ├── config/                # Configuration
│   ├── pkg/
│   │   ├── database/          # DB connection
│   │   ├── logger/            # Logging
│   │   └── middleware/        # HTTP middleware
│   └── module/workflow/       # Business logic
│   │   ├── delivery/https/    # HTTP handlers
│   │   ├── usecase/           # Business logic
│   │   ├── repository/        # Data access
│   ├── models/                # DB models
│   └── dto/                   # Request/response objects
├── database/migrations/       # Database migrations
└── configs/                   # Config files
```

---

## 📖 Swagger Integration

### Adding API Documentation

Add Swagger annotations to your handlers:

```go
// CreateSequence godoc
// @Summary Create a new email sequence
// @Description Create a sequence with multiple steps
// @Tags sequences
// @Accept json
// @Produce json
// @Param request body dto.CreateSequenceRequest true "Sequence data"
// @Success 201 {object} dto.CreateSequenceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sequences [post]
func (h *SequenceHandler) CreateSequence(c echo.Context) error {
    // Handler implementation
}
```

### Generate Documentation

```bash
make swagger
```

### View Documentation

Start server and visit: `http://localhost:8080/swagger`

---

---

## 📞 Support

* **API Docs**: <http://localhost:8080/api/v1/swagger/index.html>
* **Health Check**: <http://localhost:8080/api/v1/health>
* **Server Logs**: Check logs for debugging

---
