package https

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rohanchauhan02/sequence-service/internal/dto"
	"github.com/rohanchauhan02/sequence-service/internal/module/workflow"
	"github.com/rohanchauhan02/sequence-service/internal/pkg/ctx"
	"github.com/rohanchauhan02/sequence-service/internal/pkg/logger"
)

type workflowHandler struct {
	usecase workflow.Usecase
}

func NewWorkflowHandler(e *echo.Echo, usecase workflow.Usecase) {
	h := &workflowHandler{
		usecase: usecase,
	}

	api := e.Group("/api/v1")

	api.POST("/sequence", h.CreateSequence)
	api.PUT("/sequence/:id/steps/:stepId", h.UpdateStep)
	api.DELETE("/sequence/:id/steps/:stepId", h.DeleteStep)
	api.PATCH("/sequence/:id/tracking", h.UpdateSequence)
}

var log = logger.NewLogger()

// Handlers
func (h *workflowHandler) CreateSequence(c echo.Context) error {
	ac := c.(*ctx.CustomApplicationContext)

	reqPayload := new(dto.CreateSequenceRequest)
	if err := ac.Validate(reqPayload); err != nil {
		log.Errorf("CreateSequence - validation error: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusBadRequest), nil, "", err.Error(), http.StatusBadRequest, nil)
	}

	resp, err := h.usecase.CreateSequence(c, reqPayload)
	if err != nil {
		log.Errorf("CreateSequence - usecase error: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusInternalServerError), nil, "", err.Error(), http.StatusInternalServerError, nil)
	}

	log.Infof("CreateSequence - sequence created with ID: %s", resp.ID)
	return ac.CustomResponse("Sequence created successfully", resp, "", "", http.StatusCreated, nil)
}

func (h *workflowHandler) UpdateStep(c echo.Context) error {
	return nil
}

func (h *workflowHandler) DeleteStep(c echo.Context) error {
	return nil
}

func (h *workflowHandler) UpdateSequence(c echo.Context) error {
	return nil
}
