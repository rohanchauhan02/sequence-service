package usecase

import (
	"github.com/labstack/echo/v4"
	"github.com/rohanchauhan02/sequence-service/internal/dto"
	"github.com/rohanchauhan02/sequence-service/internal/models"
	"github.com/rohanchauhan02/sequence-service/internal/module/workflow"
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

	sequenceData := &models.Sequence{
		Name:                 req.Name,
		OpenTrackingEnabled:  req.OpenTrackingEnabled,
		ClickTrackingEnabled: req.ClickTrackingEnabled,
	}

	resp, err := u.repository.CreateSequence(nil, sequenceData)
	if err != nil {
		return nil, err
	}

	if len(req.Steps) > 0 {
		steps := make([]models.Step, len(req.Steps))
		for i, stepReq := range req.Steps {
			steps[i] = models.Step{
				SequenceID: resp.ID,
				StepOrder: stepReq.StepOrder,
				Subject:   stepReq.Subject,
				Content:   stepReq.Content,
				WaitDays:  stepReq.WaitDays,
			}
		}
		_, err = u.repository.CreateSteps(nil, steps)
		if err != nil {
			return nil, err
		}
	}

	return &dto.CreateSequenceResponse{
		ID: resp.ID.String(),
	}, nil
}
