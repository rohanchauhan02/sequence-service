package dto

type CreateSequenceRequest struct {
	Name                 string              `json:"name" validate:"required,min=1,max=255"`
	OpenTrackingEnabled  bool                `json:"open_tracking_enabled"`
	ClickTrackingEnabled bool                `json:"click_tracking_enabled"`
	Steps                []CreateStepRequest `json:"steps"`
}
type CreateSequenceResponse struct {
	ID string `json:"id"`
}
type UpdateSequenceTrackingRequest struct {
	OpenTrackingEnabled  *bool `json:"open_tracking_enabled"`
	ClickTrackingEnabled *bool `json:"click_tracking_enabled"`
}

type CreateStepRequest struct {
	StepOrder int    `json:"step_order" validate:"required,min=0"`
	Subject   string `json:"subject" validate:"required,min=1"`
	Content   string `json:"content" validate:"required,min=1"`
	WaitDays  int    `json:"wait_days" validate:"required,min=0"`
}

type UpdateStepRequest struct {
	Subject *string `json:"subject" validate:"omitempty,min=1"`
	Content *string `json:"content" validate:"omitempty,min=1"`
}

type UpdateSequenceRequest struct {
	OpenTrackingEnabled  *bool `json:"open_tracking_enabled"`
	ClickTrackingEnabled *bool `json:"click_tracking_enabled"`
}
