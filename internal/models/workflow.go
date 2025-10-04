package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Sequence struct {
	ID                   uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name                 string         `json:"name" gorm:"type:varchar(255);not null"`
	OpenTrackingEnabled  bool           `json:"open_tracking_enabled" gorm:"default:true"`
	ClickTrackingEnabled bool           `json:"click_tracking_enabled" gorm:"default:true"`
	Steps                []Step         `json:"steps,omitempty" gorm:"foreignKey:SequenceID;constraint:OnDelete:CASCADE"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
}

type Step struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SequenceID uuid.UUID      `json:"sequence_id" gorm:"type:uuid;not null;index"`
	StepOrder  int            `json:"step_order" gorm:"not null"`
	Subject    string         `json:"subject" gorm:"type:text;not null"`
	Content    string         `json:"content" gorm:"type:text;not null"`
	WaitDays   int            `json:"wait_days" gorm:"default:1"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
}
