package usecase

import "github.com/rohanchauhan02/sequence-service/internal/module/workflow"

type workflowUsecase struct {
	repository workflow.Repository
}

func NewWorkflowUsecase(repository workflow.Repository) workflow.Usecase {
	return &workflowUsecase{
		repository: repository,
	}
}
