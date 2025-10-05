package repository

import (
	"github.com/rohanchauhan02/sequence-service/internal/module/scheduler"
	"gorm.io/gorm"
)

type schedulerRepository struct {
	db *gorm.DB
}

func NewSchedulerRepository(db *gorm.DB) scheduler.Repository {
	return &schedulerRepository{
		db: db,
	}
}
