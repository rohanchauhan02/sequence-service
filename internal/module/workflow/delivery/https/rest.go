package https

import (
	"github.com/labstack/echo/v4"
	"github.com/rohanchauhan02/sequence-service/internal/module/workflow"
)

type workflowHandler struct {
	usecase workflow.Usecase
}

func NewWorkflowHandler(e *echo.Echo, usecase workflow.Usecase) {
	h := &workflowHandler{
		usecase: usecase,
	}

	api := e.Group("/api/v1/workflow")

	api.POST("/sequence", h.CreateSequence)
	api.PUT("/sequence/:id/steps/:stepId", h.UpdateStep)
	api.DELETE("/sequence/:id/steps/:stepId", h.DeleteStep)
	api.PATCH("/sequence/:id/tracking", h.UpdateTracking)
}

// Handlers

func (h *workflowHandler) CreateSequence(c echo.Context) error {
	return nil
}

func (h *workflowHandler) UpdateStep(c echo.Context) error {
	return nil
}

func (h *workflowHandler) DeleteStep(c echo.Context) error {
	return nil
}

func (h *workflowHandler) UpdateTracking(c echo.Context) error {
	return nil
}
