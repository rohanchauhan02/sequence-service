package workflow

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rohanchauhan02/sequence-service/internal/dto"
	"github.com/rohanchauhan02/sequence-service/internal/models"
	"gorm.io/gorm"
)

type Usecase interface {
	CreateSequence(c echo.Context, req *dto.CreateSequenceRequest) (*dto.CreateSequenceResponse, error)
	GetSequence(c echo.Context, sequenceID uuid.UUID) (*models.Sequence, error)
	UpdateSequenceTracking(c echo.Context, sequenceID uuid.UUID, req *dto.UpdateSequenceTrackingRequest) error
	UpdateStep(c echo.Context, sequenceID uuid.UUID, stepID uuid.UUID, req *dto.UpdateStepRequest) error
	DeleteStep(c echo.Context, sequenceID uuid.UUID, stepID uuid.UUID) error
}

type Repository interface {
	CreateSequence(tx *gorm.DB, sequence *models.Sequence) (*models.Sequence, error)
	GetSequence(sequenceID uuid.UUID) (*models.Sequence, error)
	UpdateSequenceTracking(tx *gorm.DB, sequence *models.Sequence) error

	CreateSteps(tx *gorm.DB, steps []models.Step) (*[]models.Step, error)
	GetStepByID(sequenceID uuid.UUID, stepID uuid.UUID) (*models.Step, error)

	UpdateStep(tx *gorm.DB, sequence *models.Step) error
	DeleteStep(tx *gorm.DB, sequenceID uuid.UUID, stepID uuid.UUID) error
}
