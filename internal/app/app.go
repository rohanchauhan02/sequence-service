package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rohanchauhan02/sequence-service/docs/swagger"
	EchoSwagger "github.com/swaggo/echo-swagger"

	"github.com/rohanchauhan02/sequence-service/internal/pkg/logger"
	CustomMiddleware "github.com/rohanchauhan02/sequence-service/internal/pkg/middleware"
	"github.com/rohanchauhan02/sequence-service/internal/pkg/transporter/kafka"

	"github.com/rohanchauhan02/sequence-service/internal/config"
	HealthHandler "github.com/rohanchauhan02/sequence-service/internal/module/health/delivery/https"
	HealthRepository "github.com/rohanchauhan02/sequence-service/internal/module/health/repository"
	HealthUsecase "github.com/rohanchauhan02/sequence-service/internal/module/health/usecase"
	"github.com/rohanchauhan02/sequence-service/internal/pkg/ctx"
	"github.com/rohanchauhan02/sequence-service/internal/pkg/database"
	"github.com/rohanchauhan02/sequence-service/internal/pkg/utils"

	WorkflowHandler "github.com/rohanchauhan02/sequence-service/internal/module/workflow/delivery/https"
	WorkflowRepository "github.com/rohanchauhan02/sequence-service/internal/module/workflow/repository"
	WorkflowUsecase "github.com/rohanchauhan02/sequence-service/internal/module/workflow/usecase"
)

var log = logger.NewLogger("SEQUENCE-SERVICE")

func Init() {
	e := echo.New()

	// Load configuration
	cnf := config.NewImmutableConfig()

	// Initialize database
	dbClient := database.NewPostgressClient(cnf)

	db, err := dbClient.InitClient(context.TODO())
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
		panic(err)
	}

	kafkaClient, err := kafka.NewKafkaClient(cnf)
	if err != nil {
		log.Errorf("Failed to initialize Kafka client: %v", err)
		panic(err)
	}

	defer func() {
		if err := kafkaClient.Close(); err != nil {
			log.Errorf("Failed to close Kafka client: %v", err)
		}
	}()

	// use requestID middleware
	e.Use(CustomMiddleware.MiddlewareRequestID())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())

	// Middleware to inject dependencies into the request context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			appLogger := log.WithRequestID(requestID)
			customCtx := &ctx.CustomApplicationContext{
				Context:  c,
				AppLoger: appLogger,
				Config:   cnf,
				Kakfa:    kafkaClient,
				Postgres: db,
			}
			return next(customCtx)
		}
	})

	validator := utils.DefaultValidator()
	e.Validator = validator

	setupSwaggerRoutes(e)

	// Initialize repositories
	healthRepo := HealthRepository.NewHealthRepository(db)
	workflowRepo := WorkflowRepository.NewWorkflowRepository(db)

	// Initialize usecases
	healthUsecase := HealthUsecase.NewHealthUsecase(healthRepo)
	workflowUsecase := WorkflowUsecase.NewWorkflowUsecase(workflowRepo)

	// Initialize handlers
	HealthHandler.NewHealthHandler(e, healthUsecase)
	WorkflowHandler.NewWorkflowHandler(e, workflowUsecase)

	// Start server in a separate goroutine
	serverAddr := fmt.Sprintf(":%s", cnf.GetPort())
	go func() {
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			log.Errorf("Server shutdown unexpectedly: %v", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	}
	log.Info("Server exited properly.")
}

func setupSwaggerRoutes(e *echo.Echo) {
	swagger.SwaggerInfo.Title = "Sequence Service API"
	swagger.SwaggerInfo.Description = "This is the API documentation for the Sequence Service."
	swagger.SwaggerInfo.Version = "1.0"
	swagger.SwaggerInfo.Host = "localhost:8080"
	swagger.SwaggerInfo.BasePath = "/api/v1"
	swagger.SwaggerInfo.Schemes = []string{"http", "https"}
	e.GET("/swagger/*", EchoSwagger.WrapHandler)
}
