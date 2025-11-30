package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CheckID   uuid.UUID `gorm:"type:uuid;index;not null" json:"check_id"`
	Type      string    `gorm:"type:varchar(20);not null" json:"type"`
	Payload   string    `gorm:"type:text;not null" json:"payload"`
	CreatedAt time.Time `json:"created_at"`
	Check     Check     `gorm:"foreignKey:CheckID" json:"check,omitempty"`
	SentAt    time.Time `gorm:"not null" json:"sent_at"`
}

type WebhookAttempt struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	AlertID         *uuid.UUID     `gorm:"type:uuid;index" json:"alert_id,omitempty"`
	CheckID         uuid.UUID      `gorm:"type:uuid;index;not null" json:"check_id"`
	URL             string         `gorm:"type:text;not null" json:"url"`
	RequestBody     string         `gorm:"type:text" json:"request_body"`
	RequestHeaders  datatypes.JSON `gorm:"type:jsonb" json:"request_headers"`
	ResponseCode    *int           `gorm:"type:integer" json:"response_code,omitempty"`
	ResponseBody    *string        `gorm:"type:text" json:"response_body,omitempty"`
	ResponseHeaders datatypes.JSON `gorm:"type:jsonb" json:"response_headers"`
	Error           *string        `gorm:"type:text" json:"error,omitempty"`
	LatencyMs       *int64         `gorm:"type:bigint" json:"latency_ms,omitempty"`
	RetryNumber     int            `gorm:"default:0" json:"retry_number"`
	Timeout         bool           `gorm:"default:false" json:"timeout"`
	Status          string         `gorm:"type:varchar(20);not null" json:"status"` // success, failed
	CreatedAt       time.Time      `json:"created_at"`
}
