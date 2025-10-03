package https

import (
	"github.com/labstack/echo/v4"
	"github.com/rohanchauhan02/sequence-service/internal/module/health"
)

type healthHandler struct {
	usecase health.Usecase
}

func NewHealthHandler(e *echo.Echo, usecase health.Usecase) {
	h := healthHandler{
		usecase: usecase,
	}

	api := e.Group("/api/v1")

	api.GET("/health", h.Health)
}

func (h *healthHandler) Health(c echo.Context) error {
	return c.JSON(200, "searvice is healthy")
}
