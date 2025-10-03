package repository

import (
	"github.com/rohanchauhan02/sequence-service/internal/module/workflow"
	"gorm.io/gorm"
)

type workflowRepository struct {
	db *gorm.DB
}

func NewWorkflowRepository(db *gorm.DB) workflow.Repository{
	return &workflowRepository{
		db: db,
	}
}
