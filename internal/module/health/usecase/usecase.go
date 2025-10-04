package usecase

import "github.com/rohanchauhan02/sequence-service/internal/module/health"

type healthUsecase struct {
	repository health.Repository
}

func NewHealthUsecase(repository health.Repository) health.Usecase {
	return &healthUsecase{
		repository: repository,
	}
}

func (h *healthUsecase) Health() (map[string]any, error) {
	return h.repository.Health()
}
