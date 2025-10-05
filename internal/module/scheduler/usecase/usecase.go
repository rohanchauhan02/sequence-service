package usecase

import "github.com/rohanchauhan02/sequence-service/internal/module/scheduler"

type schedulerUsecase struct {
	repo scheduler.Repository
}

func NewSchedulerUsecase(repo scheduler.Repository) scheduler.Usecase {
	return &schedulerUsecase{
		repo: repo,
	}
}
