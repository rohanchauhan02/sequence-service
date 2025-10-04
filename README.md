# Sequence Service

A high-performance email sequencing platform built with Go and PostgreSQL. Create and manage multi-step email campaigns with intelligent scheduling.

---

## 🚀 Features

- **Sequence Management**: Create multi-step email sequences
- **Contact Management**: Add contacts to sequences
- **Email Tracking**: Open and click tracking
- **RESTful API**: Clean, documented endpoints with Swagger
- **PostgreSQL**: Reliable data storage
- **Makefile**: Easy development commands

---

## 🛠 Tech Stack

- **Go 1.23+** - Backend language
- **Echo** - HTTP framework
- **PostgreSQL** - Database
- **GORM** - ORM
- **Viper** - Configuration
- **Goose** - Migrations
- **Swagger** - API documentation

---

## ⚡ Quick Start

### 1. Clone & Setup

```bash
git clone https://github.com/rohanchauhan02/sequence-service
cd sequence-service
cp configs/app.config.sample.yml configs/app.config.local.yml
```

### 2. Configure Database

Edit `configs/app.config.local.yml`:

```yaml
PORT: 8080

DB:
  HOST: localhost
  PORT: 5432
  NAME: sequence_db
  USER: your_username
  PASSWORD: your_password
  SSL_MODE: disable
```

### 3. Setup Database & Start

```bash
make migrate-up
make run
```

Server: **<http://localhost:8080>**
API Docs: **<http://localhost:8080/swagger>**

---

## 🛠 Development Commands

### Makefile Commands

```bash
# Start development
make run

# Run tests
make test

# Database migrations
make migrate-up
make migrate-down
make migrate-status

# Generate Swagger docs
make swagger

# Code quality
make lint
make fmt

# Clean build
make clean
```

### Manual Commands

```bash
# Install dependencies
go mod tidy

# Run tests with coverage
go test -cover ./...

# Generate Swagger docs
swag init -g cmd/api/main.go -o docs
```

---

## 📚 API Documentation

### Swagger UI

After starting the server, visit: **<http://localhost:8080/swagger>**

![Swagger UI](http://localhost:8080/swagger/index.html)

### API Examples

#### Health Check

```http
GET /api/v1/health
```

Response:

```json
{
  "request_id": "4616d193-556e-4e7d-9ad1-a75a44fb4c3a",
  "status": "Service is healthy",
  "code": 200
}
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

## 🗄 Database

### Main Tables

- `sequences` - Sequence definitions
- `steps` - Email steps
- `mailboxes` - Email accounts
- `contacts` - Recipients
- `sequence_contacts` - Links contacts to sequences
- `email_queues` - Scheduled emails

### Migration Commands

```bash
make migrate-up          # Run migrations
make migrate-down        # Rollback last migration
make migrate-status      # Check migration status
make migrate-create name=add_feature  # Create new migration
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

## 🔧 Development

### Code Quality

```bash
make fmt    # Format code
make lint   # Lint code
make test   # Run tests
```

### Database

```bash
make migrate-up     # Run migrations
make migrate-down   # Rollback
make migrate-status # Check status
```

### Building

```bash
make build    # Build binary
make run      # Run locally
make clean    # Clean build
```

---

## 📞 Support

- **API Docs**: <http://localhost:8080/api/v1/swagger/index.html>
- **Health Check**: <http://localhost:8080/api/v1/health>
- **GitHub Issues**: For bugs and feature requests
- **Server Logs**: Check logs for debugging

---
