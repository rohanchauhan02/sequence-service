package usecase

import (
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rohanchauhan02/sequence-service/internal/dto"
)

func Test_workflowUsecase_CreateSequence(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		c       echo.Context
		req     *dto.CreateSequenceRequest
		want    *dto.CreateSequenceResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var u workflowUsecase
			got, gotErr := u.CreateSequence(tt.c, tt.req)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CreateSequence() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CreateSequence() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("CreateSequence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_workflowUsecase_UpdateSequenceTracking(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		c          echo.Context
		sequenceID uuid.UUID
		req        *dto.UpdateSequenceTrackingRequest
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var u workflowUsecase
			gotErr := u.UpdateSequenceTracking(tt.c, tt.sequenceID, tt.req)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("UpdateSequenceTracking() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("UpdateSequenceTracking() succeeded unexpectedly")
			}
		})
	}
}

func Test_workflowUsecase_UpdateStep(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		c          echo.Context
		sequenceID uuid.UUID
		stepID     uuid.UUID
		req        *dto.UpdateStepRequest
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var u workflowUsecase
			gotErr := u.UpdateStep(tt.c, tt.sequenceID, tt.stepID, tt.req)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("UpdateStep() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("UpdateStep() succeeded unexpectedly")
			}
		})
	}
}

func Test_workflowUsecase_DeleteStep(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		c          echo.Context
		sequenceID uuid.UUID
		stepID     uuid.UUID
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var u workflowUsecase
			gotErr := u.DeleteStep(tt.c, tt.sequenceID, tt.stepID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("DeleteStep() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("DeleteStep() succeeded unexpectedly")
			}
		})
	}
}
