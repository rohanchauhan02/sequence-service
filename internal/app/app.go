package app

import (
	"github.com/labstack/echo/v4"

	HealthHandler "github.com/rohanchauhan02/sequence-service/internal/module/health/delivery/https"
	HealthRepository "github.com/rohanchauhan02/sequence-service/internal/module/health/repository"
	HealthUsecase "github.com/rohanchauhan02/sequence-service/internal/module/health/usecase"
)

func Init() {
	e := echo.New()

	healthRepo := HealthRepository.NewHealthRepository(nil)

	healthUsecase := HealthUsecase.NewHealthUsecase(healthRepo)

	HealthHandler.NewHealthHandler(e, healthUsecase)

	if err := e.Start(":8001"); err != nil {
		panic(err)
	}
}
