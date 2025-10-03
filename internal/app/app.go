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
	"github.com/labstack/gommon/log"

	"github.com/rohanchauhan02/sequence-service/internal/config"
	HealthHandler "github.com/rohanchauhan02/sequence-service/internal/module/health/delivery/https"
	HealthRepository "github.com/rohanchauhan02/sequence-service/internal/module/health/repository"
	HealthUsecase "github.com/rohanchauhan02/sequence-service/internal/module/health/usecase"
	"github.com/rohanchauhan02/sequence-service/internal/pkg/database"

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
