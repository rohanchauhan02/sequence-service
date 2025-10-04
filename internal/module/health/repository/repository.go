package repository

import (
	"github.com/rohanchauhan02/sequence-service/internal/module/health"
	"gorm.io/gorm"
)

type healthRepository struct {
	db *gorm.DB
}

func NewHealthRepository(db *gorm.DB) health.Repository {
	return &healthRepository{
		db: db,
	}
}

func (r *healthRepository) Health() (map[string]any, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}
	err = sqlDB.Ping()
	if err != nil {
		return map[string]any{"status": "unhealthy"}, nil
	}
	return map[string]any{"status": "healthy"}, nil
}
