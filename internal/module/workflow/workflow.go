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
}

type Repository interface {
	CreateSequence(tx *gorm.DB, sequence *models.Sequence) (*models.Sequence, error)
	CreateSteps(tx *gorm.DB, steps []models.Step) (*[]models.Step, error)
	UpdateStep(tx *gorm.DB, sequence *models.Step) error
	DeleteStep(tx *gorm.DB, id uuid.UUID) error
}
