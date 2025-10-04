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

func (r *workflowRepository) CreateSequence(tx *gorm.DB, sequence *models.Sequence) (*models.Sequence, error) {
	if err := tx.Create(sequence).Error; err != nil {
		return nil, err
	}
	return sequence, nil
}

func (r *workflowRepository) GetSequence(sequenceID uuid.UUID) (*models.Sequence, error) {
	var sequence models.Sequence
	if err := r.db.Preload("Steps", func(db *gorm.DB) *gorm.DB {
		return db.Order("step_order ASC")
	}).First(&sequence, "id = ?", sequenceID).Error; err != nil {
		return nil, err
	}
	return &sequence, nil
}

func (r *workflowRepository) UpdateSequenceTracking(tx *gorm.DB, sequence *models.Sequence) error {
	return tx.Save(sequence).Error
}

func (r *workflowRepository) CreateSteps(tx *gorm.DB, steps []models.Step) (*[]models.Step, error) {
	if err := tx.Create(&steps).Error; err != nil {
		return nil, err
	}
	return &steps, nil
}

func (r *workflowRepository) GetStepByID(sequenceID, stepID uuid.UUID) (*models.Step, error) {
	var step models.Step
	if err := r.db.Where("id = ? AND sequence_id = ?", stepID, sequenceID).First(&step).Error; err != nil {
		return nil, err
	}
	return &step, nil
}

func (r *workflowRepository) UpdateStep(tx *gorm.DB, step *models.Step) error {
	return tx.Save(step).Error
}

func (r *workflowRepository) DeleteStep(tx *gorm.DB, sequenceID, stepID uuid.UUID) error {
	return tx.Delete(&models.Step{}, "id = ? AND sequence_id = ?", stepID, sequenceID).Error
}
