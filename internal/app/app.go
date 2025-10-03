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
	"github.com/labstack/gommon/log"

	CustomMiddleware "github.com/rohanchauhan02/sequence-service/internal/pkg/middleware"

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

func Init() {
	e := echo.New()

	// Load configuration
	cnf := config.NewImmutableConfig()

	// Initialize database
	dbClient := database.NewPostgressClient(cnf)

	db, err := dbClient.InitClient(context.TODO())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		panic(err)
	}

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
			customCtx := &ctx.CustomApplicationContext{
				Context:    c,
				Config:     cnf,
				PostgresDB: db,
			}
			return next(customCtx)
		}
	})
	e.Use(middleware.CORS())

	validator := utils.DefaultValidator()
	e.Validator = validator

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
			log.Fatalf("Server shutdown unexpectedly: %v", err)
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
