package usecase

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rohanchauhan02/sequence-service/internal/dto"
	"github.com/rohanchauhan02/sequence-service/internal/models"
	"github.com/rohanchauhan02/sequence-service/internal/module/workflow"
	"github.com/rohanchauhan02/sequence-service/internal/pkg/ctx"
)

type workflowUsecase struct {
	repository workflow.Repository
}

func NewWorkflowUsecase(repository workflow.Repository) workflow.Usecase {
	return &workflowUsecase{
		repository: repository,
	}
}

func (u *workflowUsecase) CreateSequence(c echo.Context, req *dto.CreateSequenceRequest) (*dto.CreateSequenceResponse, error) {
	ac := c.(*ctx.CustomApplicationContext)
	sequenceData := &models.Sequence{
		Name:                 req.Name,
		OpenTrackingEnabled:  req.OpenTrackingEnabled,
		ClickTrackingEnabled: req.ClickTrackingEnabled,
	}

	tx := ac.PostgresDB.Begin()
	defer tx.Rollback()

	resp, err := u.repository.CreateSequence(tx, sequenceData)
	if err != nil {
		ac.AppLoger.Errorf("CreateSequence - failed to create sequence: %v", err)
		return nil, err
	}

	if len(req.Steps) > 0 {
		steps := make([]models.Step, len(req.Steps))
		for i, stepReq := range req.Steps {
			steps[i] = models.Step{
				SequenceID: resp.ID,
				StepOrder:  stepReq.StepOrder,
				Subject:    stepReq.Subject,
				Content:    stepReq.Content,
				WaitDays:   stepReq.WaitDays,
			}
		}
		_, err = u.repository.CreateSteps(tx, steps)
		if err != nil {
			ac.AppLoger.Errorf("CreateSequence - failed to create steps: %v", err)
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		ac.AppLoger.Errorf("CreateSequence - failed to commit transaction: %v", err)
		return nil, err
	}

	return &dto.CreateSequenceResponse{
		ID: resp.ID.String(),
	}, nil
}

func (u *workflowUsecase) GetSequence(c echo.Context, sequenceID uuid.UUID) (*models.Sequence, error) {
	return u.repository.GetSequence(sequenceID)
}

func (u *workflowUsecase) UpdateStep(c echo.Context, sequenceID uuid.UUID, stepID uuid.UUID, req *dto.UpdateStepRequest) error {
	ac := c.(*ctx.CustomApplicationContext)
	existingStep, err := u.repository.GetStepByID(sequenceID, stepID)
	if existingStep == nil {
		return echo.NewHTTPError(404, "Step not found")
	}
	if err != nil {
		return err
	}

	if req.Subject != nil {
		existingStep.Subject = *req.Subject
	}

	if req.Content != nil {
		existingStep.Content = *req.Content
	}

	tx := ac.PostgresDB.Begin()
	defer tx.Rollback()

	err = u.repository.UpdateStep(tx, existingStep)
	if err != nil {
		ac.AppLoger.Errorf("UpdateStep - failed to update step: %v", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		ac.AppLoger.Errorf("UpdateStep - failed to commit transaction: %v", err)
		return err
	}

	return nil
}

func (u *workflowUsecase) DeleteStep(c echo.Context, sequenceID uuid.UUID, stepID uuid.UUID) error {
	ac := c.(*ctx.CustomApplicationContext)
	tx := ac.PostgresDB.Begin()
	defer tx.Rollback()

	err := u.repository.DeleteStep(nil, sequenceID, stepID)
	if err != nil {
		ac.AppLoger.Errorf("DeleteStep - failed to delete step: %v", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		ac.AppLoger.Errorf("DeleteStep - failed to commit transaction: %v", err)
		return err
	}

	return nil
}
