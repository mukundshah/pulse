package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type Check struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name             string         `gorm:"not null" json:"name"`
	URL              string         `gorm:"not null" json:"url"`
	Method           string         `gorm:"not null;default:GET" json:"method"`
	Headers          datatypes.JSON `gorm:"type:jsonb" json:"headers"`
	ExpectedStatus   int            `gorm:"default:200" json:"expected_status"`
	BodyContains     *string        `gorm:"type:text" json:"body_contains,omitempty"`
	TimeoutMs        int            `gorm:"default:10000" json:"timeout_ms"`        // timeout in milliseconds
	WebhookURL       *string        `gorm:"type:text" json:"webhook_url,omitempty"` // webhook URL for alerts
	IntervalSeconds  int            `gorm:"not null" json:"interval_seconds"`
	AlertThreshold   int            `gorm:"not null;default:3" json:"alert_threshold"`
	ConsecutiveFails int            `gorm:"not null;default:0" json:"consecutive_fails"`
	LastStatus       string         `gorm:"type:varchar(20);default:'unknown'" json:"last_status"`
	LastRunAt        *time.Time     `json:"last_run_at,omitempty"`
	NextRunAt        *time.Time     `json:"next_run_at,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

type CheckRun struct {
	ID         uuid.UUID `json:"id"`
	CheckID    uuid.UUID `json:"check_id"`
	Status     string    `json:"status"` // success, fail, timeout, error
	LatencyMs  int64     `json:"latency_ms"`
	StatusCode int32     `json:"status_code"`
	Error      *string   `json:"error,omitempty"`
	RunAt      time.Time `json:"run_at"`
}

type Alert struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CheckID   uuid.UUID `gorm:"type:uuid;index;not null"`
	Type      string    `gorm:"type:varchar(20);not null"`
	Payload   string    `gorm:"type:text;not null"`
	CreatedAt time.Time
	Check     Check     `gorm:"foreignKey:CheckID"`
	SentAt    time.Time `gorm:"not null"`
}

type WebhookAttempt struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey"`
	AlertID        *uuid.UUID     `gorm:"type:uuid;index"`
	CheckID        uuid.UUID      `gorm:"type:uuid;index;not null"`
	URL            string         `gorm:"type:text;not null"`
	RequestBody    string         `gorm:"type:text"`
	RequestHeaders datatypes.JSON `gorm:"type:jsonb"`
	ResponseCode   *int           `gorm:"type:integer"`
	ResponseBody   *string         `gorm:"type:text"`
	ResponseHeaders datatypes.JSON `gorm:"type:jsonb"`
	Error          *string         `gorm:"type:text"`
	LatencyMs      *int64         `gorm:"type:bigint"`
	RetryNumber    int            `gorm:"default:0"`
	Timeout        bool           `gorm:"default:false"`
	Status         string         `gorm:"type:varchar(20);not null"` // success, failed
	CreatedAt      time.Time
}
