package repository

import (
	"github.com/rohanchauhan02/sequence-service/internal/module/health"
	"gorm.io/gorm"
)

type healthRepository struct {
	db *gorm.DB
}

func NewHealthRepository(db *gorm.DB) health.Repository{
	return &healthRepository{
		db: db,
	}
}
