package repository

import (
	"github.com/google/uuid"
	"github.com/rohanchauhan02/sequence-service/internal/models"
	"github.com/rohanchauhan02/sequence-service/internal/module/workflow"
	"gorm.io/gorm"
)

type workflowRepository struct {
	db *gorm.DB
}

func NewWorkflowRepository(db *gorm.DB) workflow.Repository {
	return &workflowRepository{
		db: db,
	}
}

func (r *workflowRepository) CreateSequence(tx *gorm.DB, sequence *models.Sequence) (*models.Sequence, error){
	if err := r.db.Create(sequence).Error;  err != nil {
		return nil, err
	}
	return sequence, nil
}

func (r *workflowRepository) CreateSteps(tx *gorm.DB, steps []models.Step) (*[]models.Step, error) {
	if err := r.db.Create(&steps).Error; err != nil {
		return nil, err
	}
	return &steps, nil
}

func (r *workflowRepository) UpdateStep(tx *gorm.DB, step *models.Step) error {
	return tx.Save(step).Error
}

func (r *workflowRepository) DeleteStep(tx *gorm.DB, id uuid.UUID) error {
	return tx.Delete(&models.Step{}, "id = ?", id).Error
}
