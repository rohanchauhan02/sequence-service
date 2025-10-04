package https

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rohanchauhan02/sequence-service/internal/dto"
	"github.com/rohanchauhan02/sequence-service/internal/module/workflow"
	"github.com/rohanchauhan02/sequence-service/internal/pkg/ctx"
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
	api.PATCH("/sequence/:id", h.UpdateSequenceTracking)
}

// CreateSequence godoc
// @Summary      Create a new email sequence
// @Description  Create a new email sequence with steps
// @Tags         Sequences
// @Accept       json
// @Produce      json
// @Param        sequence  body      dto.CreateSequenceRequest  true  "Sequence details"
// @Success      201  {object}  dto.CreateSequenceResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /sequence [post]
func (h *workflowHandler) CreateSequence(c echo.Context) error {
	ac := c.(*ctx.CustomApplicationContext)

	reqPayload := new(dto.CreateSequenceRequest)
	if err := ac.CustomBind(reqPayload); err != nil {
		ac.AppLoger.Errorf("CreateSequence - validation error: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusBadRequest), nil, "", err.Error(), http.StatusBadRequest, nil)
	}

	resp, err := h.usecase.CreateSequence(c, reqPayload)
	if err != nil {
		ac.AppLoger.Errorf("CreateSequence - usecase error: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusInternalServerError), nil, "", err.Error(), http.StatusInternalServerError, nil)
	}

	ac.AppLoger.Infof("CreateSequence - sequence created with ID: %s", resp.ID)
	return ac.CustomResponse("Sequence created successfully", resp, "", "", http.StatusCreated, nil)
}

// GetSequence godoc
// @Summary      Get sequence details
// @Description  Retrieve details of a specific email sequence by its ID
// @Tags         Sequences
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Sequence ID"
// @Success      200  {object}  models.Sequence
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /sequence/{id} [get]
func (h *workflowHandler) GetSequence(c echo.Context) error {
	ac := c.(*ctx.CustomApplicationContext)

	sequenceID := c.Param("id")
	ac.AppLoger.Infof("GetSequence - sequenceID: %s", sequenceID)
	sequenceUUID, err := uuid.Parse(sequenceID)
	if err != nil {
		ac.AppLoger.Errorf("GetSequence - invalid sequence ID: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusBadRequest), nil, "", "Invalid sequence ID", http.StatusBadRequest, nil)
	}

	sequenceDetails, err := h.usecase.GetSequence(c, sequenceUUID)
	if err != nil {
		ac.AppLoger.Errorf("GetSequence - usecase error: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusInternalServerError), nil, "", err.Error(), http.StatusInternalServerError, nil)
	}
	ac.AppLoger.Infof("GetSequence - sequence details retrieved for ID: %s", sequenceID)
	return ac.CustomResponse("Sequence details retrieved successfully", sequenceDetails, "", "", http.StatusOK, nil)
}

// UpdateStep godoc
// @Summary      Update a step in the sequence
// @Description  Update details of a specific step within an email sequence
// @Tags         Sequences
// @Accept       json
// @Produce      json
// @Param        id      path      string                   true  "Sequence ID"
// @Param        stepId  path      string                   true  "Step ID"
// @Param        step    body      dto.UpdateStepRequest    true  "Step details to update"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /sequence/{id}/steps/{stepId} [put]
func (h *workflowHandler) UpdateStep(c echo.Context) error {
	ac := c.(*ctx.CustomApplicationContext)

	reqPayload := new(dto.UpdateStepRequest)
	if err := ac.CustomBind(reqPayload); err != nil {
		ac.AppLoger.Errorf("UpdateStep - validation error: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusBadRequest), nil, "", err.Error(), http.StatusBadRequest, nil)
	}

	sequenceID := c.Param("id")
	stepID := c.Param("stepId")

	ac.AppLoger.Infof("UpdateStep - sequenceID: %s, stepID: %s, payload: %+v", sequenceID, stepID, reqPayload)

	sequenceUUID, err := uuid.Parse(sequenceID)
	if err != nil {
		ac.AppLoger.Errorf("UpdateStep - invalid sequence ID: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusBadRequest), nil, "", "Invalid sequence ID", http.StatusBadRequest, nil)
	}

	stepUUID, err := uuid.Parse(stepID)
	if err != nil {
		ac.AppLoger.Errorf("UpdateStep - invalid step ID: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusBadRequest), nil, "", "Invalid step ID", http.StatusBadRequest, nil)
	}

	err = h.usecase.UpdateStep(c, sequenceUUID, stepUUID, reqPayload)
	if err != nil {
		ac.AppLoger.Errorf("UpdateStep - usecase error: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusInternalServerError), nil, "", err.Error(), http.StatusInternalServerError, nil)
	}

	ac.AppLoger.Infof("UpdateStep - step updated successfully for sequenceID: %s, stepID: %s", sequenceID, stepID)
	return ac.CustomResponse("Step updated successfully", nil, "", "", http.StatusOK, nil)
}

// DeleteStep godoc
// @Summary      Delete a step from the sequence
// @Description  Remove a specific step from an email sequence
// @Tags         Sequences
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Sequence ID"
// @Param        stepId  path      string  true  "Step ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /sequence/{id}/steps/{stepId} [delete]
func (h *workflowHandler) DeleteStep(c echo.Context) error {
	ac := c.(*ctx.CustomApplicationContext)

	sequenceID := c.Param("id")
	stepID := c.Param("stepId")

	ac.AppLoger.Infof("DeleteStep - sequenceID: %s, stepID: %s", sequenceID, stepID)
	sequenceUUID, err := uuid.Parse(sequenceID)
	if err != nil {
		ac.AppLoger.Errorf("DeleteStep - invalid sequence ID: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusBadRequest), nil, "", "Invalid sequence ID", http.StatusBadRequest, nil)
	}

	stepUUID, err := uuid.Parse(stepID)
	if err != nil {
		ac.AppLoger.Errorf("DeleteStep - invalid step ID: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusBadRequest), nil, "", "Invalid step ID", http.StatusBadRequest, nil)
	}

	err = h.usecase.DeleteStep(c, sequenceUUID, stepUUID)
	if err != nil {
		ac.AppLoger.Errorf("DeleteStep - usecase error: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusInternalServerError), nil, "", err.Error(), http.StatusInternalServerError, nil)
	}

	ac.AppLoger.Infof("DeleteStep - step deleted successfully for sequenceID: %s, stepID: %s", sequenceID, stepID)
	return ac.CustomResponse("Step deleted successfully", nil, "", "", http.StatusOK, nil)
}

// UpdateSequenceTracking godoc
// @Summary      Update sequence tracking information
// @Description  Update tracking information for a specific email sequence
// @Tags         Sequences
// @Accept       json
// @Produce      json
// @Param        id        path      string                           true  "Sequence ID"
// @Param        tracking  body      dto.UpdateSequenceTrackingRequest  true  "Tracking information to update"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /sequence/{id} [patch]
func (h *workflowHandler) UpdateSequenceTracking(c echo.Context) error {

	ac := c.(*ctx.CustomApplicationContext)

	sequenceID := c.Param("id")
	ac.AppLoger.Infof("UpdateSequenceTracking - sequenceID: %s", sequenceID)
	sequenceUUID, err := uuid.Parse(sequenceID)
	if err != nil {
		ac.AppLoger.Errorf("UpdateSequenceTracking - invalid sequence ID: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusBadRequest), nil, "", "Invalid sequence ID", http.StatusBadRequest, nil)
	}

	reqPayload := new(dto.UpdateSequenceTrackingRequest)
	if err := ac.CustomBind(reqPayload); err != nil {
		ac.AppLoger.Errorf("UpdateSequenceTracking - validation error: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusBadRequest), nil, "", err.Error(), http.StatusBadRequest, nil)
	}

	err = h.usecase.UpdateSequenceTracking(c, sequenceUUID, reqPayload)
	if err != nil {
		ac.AppLoger.Errorf("UpdateSequenceTracking - usecase error: %v", err)
		return ac.CustomResponse(http.StatusText(http.StatusInternalServerError), nil, "", err.Error(), http.StatusInternalServerError, nil)
	}

	ac.AppLoger.Infof("UpdateSequenceTracking - tracking info updated for sequence ID: %s", sequenceID)

	return ac.CustomResponse("Sequence tracking info updated successfully", map[string]string{"sequence_id": sequenceUUID.String()}, "", "", http.StatusOK, nil)
}
