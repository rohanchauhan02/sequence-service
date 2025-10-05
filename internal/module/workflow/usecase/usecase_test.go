package usecase

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	mock_workflow "github.com/rohanchauhan02/sequence-service/files/mocks/workflow"
	"github.com/rohanchauhan02/sequence-service/internal/dto"
	"github.com/rohanchauhan02/sequence-service/internal/models"
	"github.com/rohanchauhan02/sequence-service/internal/pkg/ctx"
	"github.com/rohanchauhan02/sequence-service/internal/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockCtx(db *gorm.DB) echo.Context {
	e := echo.New()
	return &ctx.CustomApplicationContext{
		Context:    e.NewContext(nil, nil),
		Postgres: db,
		AppLoger:   logger.NewLogger(),
	}
}

func Test_CreateSequence(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_workflow.NewMockRepository(ctrl)
	u := NewWorkflowUsecase(mockRepo)

	tests := []struct {
		name       string
		req        *dto.CreateSequenceRequest
		setupMocks func()
		wantErr    bool
	}{
		{
			name: "success - create sequence with steps",
			req: &dto.CreateSequenceRequest{
				Name:                 "Seq 1",
				OpenTrackingEnabled:  true,
				ClickTrackingEnabled: true,
				Steps: []dto.CreateStepRequest{
					{StepOrder: 1, Subject: "Subj", Content: "Cont", WaitDays: 2},
				},
			},
			setupMocks: func() {
				mockRepo.EXPECT().
					CreateSequence(gomock.Any(), gomock.Any()).
					Return(&models.Sequence{ID: uuid.New()}, nil)

				mockRepo.EXPECT().
					CreateSteps(gomock.Any(), gomock.Any()).
					Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name: "error - repository fails to create sequence",
			req:  &dto.CreateSequenceRequest{Name: "Bad Seq"},
			setupMocks: func() {
				mockRepo.EXPECT().
					CreateSequence(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			if !tt.wantErr {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			c := newMockCtx(gormDB)

			tt.setupMocks()

			_, err := u.CreateSequence(c, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSequence() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet SQL expectations: %v", err)
			}
		})
	}

}
