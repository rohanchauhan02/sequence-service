package app

import (
	"github.com/labstack/echo/v4"

	HealthHandler "github.com/rohanchauhan02/sequence-service/internal/module/health/delivery/https"
	HealthRepository "github.com/rohanchauhan02/sequence-service/internal/module/health/repository"
	HealthUsecase "github.com/rohanchauhan02/sequence-service/internal/module/health/usecase"

	WorkflowHandler "github.com/rohanchauhan02/sequence-service/internal/module/workflow/delivery/https"
	WorkflowRepository "github.com/rohanchauhan02/sequence-service/internal/module/workflow/repository"
	WorkflowUsecase "github.com/rohanchauhan02/sequence-service/internal/module/workflow/usecase"
)

func Init() {
	e := echo.New()

	healthRepo := HealthRepository.NewHealthRepository(nil)
	healthUsecase := HealthUsecase.NewHealthUsecase(healthRepo)
	HealthHandler.NewHealthHandler(e, healthUsecase)

	workflowRepo := WorkflowRepository.NewWorkflowRepository(nil)
	workflowUsecase := WorkflowUsecase.NewWorkflowUsecase(workflowRepo)
	WorkflowHandler.NewWorkflowHandler(e, workflowUsecase)

	if err := e.Start(":8001"); err != nil {
		panic(err)
	}
}
