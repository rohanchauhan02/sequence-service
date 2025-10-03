package https

import (
	"net/http"

	"github.com/google/uuid"
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
	api.GET("/sequence/:id", h.GetSequence)
	api.PUT("/sequence/:id/steps/:stepId", h.UpdateStep)
	api.DELETE("/sequence/:id/steps/:stepId", h.DeleteStep)
	api.PATCH("/sequence/:id/tracking", h.UpdateSequence)
}

var log = logger.NewLogger()

// Handlers
func (h *workflowHandler) CreateSequence(c echo.Context) error {
	ac := c.(*ctx.CustomApplicationContext)

	reqPayload := new(dto.CreateSequenceRequest)
	if err := ac.CustomBind(reqPayload); err != nil {
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

func (h *workflowHandler) GetSequence(c echo.Context) error {
	ac := c.(*ctx.CustomApplicationContext)

	sequenceID := c.Param("id")
	log.Infof("GetSequence - sequenceID: %s", sequenceID)
	sequenceUUID, err := uuid.Parse(sequenceID)
	if err != nil {
		log.Errorf("GetSequence - invalid sequence ID: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusBadRequest), nil, "", "Invalid sequence ID", http.StatusBadRequest, nil)
	}

	sequenceDetails, err := h.usecase.GetSequence(c, sequenceUUID)
	if err != nil {
		log.Errorf("GetSequence - usecase error: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusInternalServerError), nil, "", err.Error(), http.StatusInternalServerError, nil)
	}
	log.Infof("GetSequence - sequence details retrieved for ID: %s", sequenceID)
	return ac.CustomResponse("Sequence details retrieved successfully", sequenceDetails, "", "", http.StatusOK, nil)
}

func (h *workflowHandler) UpdateStep(c echo.Context) error {
	ac := c.(*ctx.CustomApplicationContext)

	reqPayload := new(dto.UpdateStepRequest)
	if err := ac.CustomBind(reqPayload); err != nil {
		log.Errorf("UpdateStep - validation error: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusBadRequest), nil, "", err.Error(), http.StatusBadRequest, nil)
	}

	sequenceID := c.Param("id")
	stepID := c.Param("stepId")

	log.Infof("UpdateStep - sequenceID: %s, stepID: %s, payload: %+v", sequenceID, stepID, reqPayload)

	sequenceUUID, err := uuid.Parse(sequenceID)
	if err != nil {
		log.Errorf("UpdateStep - invalid sequence ID: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusBadRequest), nil, "", "Invalid sequence ID", http.StatusBadRequest, nil)
	}

	stepUUID, err := uuid.Parse(stepID)
	if err != nil {
		log.Errorf("UpdateStep - invalid step ID: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusBadRequest), nil, "", "Invalid step ID", http.StatusBadRequest, nil)
	}

	err = h.usecase.UpdateStep(c, sequenceUUID, stepUUID, reqPayload)
	if err != nil {
		log.Errorf("UpdateStep - usecase error: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusInternalServerError), nil, "", err.Error(), http.StatusInternalServerError, nil)
	}

	log.Infof("UpdateStep - step updated successfully for sequenceID: %s, stepID: %s", sequenceID, stepID)
	return ac.CustomResponse("Step updated successfully", nil, "", "", http.StatusOK, nil)
}

func (h *workflowHandler) DeleteStep(c echo.Context) error {
	ac := c.(*ctx.CustomApplicationContext)

	sequenceID := c.Param("id")
	stepID := c.Param("stepId")

	log.Infof("DeleteStep - sequenceID: %s, stepID: %s", sequenceID, stepID)
	sequenceUUID, err := uuid.Parse(sequenceID)
	if err != nil {
		log.Errorf("DeleteStep - invalid sequence ID: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusBadRequest), nil, "", "Invalid sequence ID", http.StatusBadRequest, nil)
	}

	stepUUID, err := uuid.Parse(stepID)
	if err != nil {
		log.Errorf("DeleteStep - invalid step ID: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusBadRequest), nil, "", "Invalid step ID", http.StatusBadRequest, nil)
	}

	err = h.usecase.DeleteStep(c, sequenceUUID, stepUUID)
	if err != nil {
		log.Errorf("DeleteStep - usecase error: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusInternalServerError), nil, "", err.Error(), http.StatusInternalServerError, nil)
	}

	log.Infof("DeleteStep - step deleted successfully for sequenceID: %s, stepID: %s", sequenceID, stepID)
	return ac.CustomResponse("Step deleted successfully", nil, "", "", http.StatusOK, nil)
}

func (h *workflowHandler) UpdateSequence(c echo.Context) error {
	return nil
}
