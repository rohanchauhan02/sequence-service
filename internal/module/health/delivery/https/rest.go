package https

import (
	"github.com/labstack/echo/v4"
	"github.com/rohanchauhan02/sequence-service/internal/module/health"
	"github.com/rohanchauhan02/sequence-service/internal/pkg/ctx"
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

// Health godoc
// @Summary      Check the health status of the service
// @Description  Returns the health status of the service
// @Tags         Health
// @Produce      json
// @Success      200  {object}  dto.ResponsePattern
// @Failure      500  {object}  dto.ResponsePattern
// @Router       /health [get]
func (h *healthHandler) Health(c echo.Context) error {
	ac := c.(*ctx.CustomApplicationContext)

	resp, err := h.usecase.Health()
	if err != nil {
		ac.AppLoger.Errorf("Health - usecase error: %v", err)
		return ac.CustomResponse("Service is unhealthy", nil, "", err.Error(), 500, nil)
	}

	return ac.CustomResponse("Service is healthy", resp, "Service is healthy", "", 200, nil)
}
