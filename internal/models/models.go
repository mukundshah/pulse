package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CheckType string

const (
	CheckTypeHTTP      CheckType = "http"
	CheckTypeTCP       CheckType = "tcp"
	CheckTypeDNS       CheckType = "dns"
	CheckTypeBrowser   CheckType = "browser"
	CheckTypeHeartbeat CheckType = "heartbeat"
)

type CheckRunStatus string

const (
	CheckRunStatusSuccess CheckRunStatus = "success"
	CheckRunStatusFail    CheckRunStatus = "fail"
	CheckRunStatusTimeout CheckRunStatus = "timeout"
	CheckRunStatusError   CheckRunStatus = "error"
	CheckRunStatusUnknown CheckRunStatus = "unknown"
)

type Tag struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Project struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Tags []Tag `gorm:"many2many:project_tags;" json:"tags,omitempty"`
}

type Region struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Code      string         `gorm:"not null;unique" json:"code"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Check struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	Name string    `gorm:"not null" json:"name"`

	IsEnabled bool `gorm:"default:true" json:"is_enabled"`
	IsMuted   bool `gorm:"default:false" json:"is_muted"`

	Type CheckType `gorm:"not null;default:http" json:"type"`

	URL     string         `gorm:"not null" json:"url"`
	Method  string         `gorm:"not null;default:GET" json:"method"`
	Headers datatypes.JSON `gorm:"type:jsonb" json:"headers"`

	PlaywrightScript *string        `gorm:"type:text" json:"playwright_script,omitempty"`
	Assertions       datatypes.JSON `gorm:"type:jsonb" json:"assertions"`
	ExpectedStatus   int            `gorm:"default:200" json:"expected_status"`
	ShouldFail       bool           `gorm:"default:false" json:"should_fail"`

	PreScript  *string `gorm:"type:text" json:"pre_script,omitempty"`
	PostScript *string `gorm:"type:text" json:"post_script,omitempty"`

	TimeoutMs        int `gorm:"default:10000" json:"timeout_ms"`
	IntervalSeconds  int `gorm:"not null" json:"interval_seconds"`
	AlertThreshold   int `gorm:"not null;default:3" json:"alert_threshold"`
	ConsecutiveFails int `gorm:"not null;default:0" json:"consecutive_fails"`

	LastStatus CheckRunStatus `gorm:"type:varchar(20);default:'unknown'" json:"last_status"`
	LastRunAt  *time.Time     `json:"last_run_at,omitempty"`
	NextRunAt  *time.Time     `json:"next_run_at,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ProjectID uuid.UUID `gorm:"type:uuid;index;not null" json:"project_id"`
	Project   Project   `gorm:"foreignKey:ProjectID" json:"project,omitempty"`

	Tags    []Tag    `gorm:"many2many:check_tags;" json:"tags,omitempty"`
	Regions []Region `gorm:"many2many:check_regions;" json:"regions,omitempty"`
}

type CheckRun struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`

	Status         CheckRunStatus `gorm:"type:varchar(20);default:'unknown'" json:"status"`
	ResponseStatus int32          `gorm:"type:integer" json:"response_status"`

	AssertionResults datatypes.JSON `gorm:"type:jsonb" json:"assertion_results"`
	PlaywrightReport datatypes.JSON `gorm:"type:jsonb" json:"playwright_report,omitempty"`
	NetworkTimings   datatypes.JSON `gorm:"type:jsonb" json:"network_timings"`
	Metrics          datatypes.JSON `gorm:"type:jsonb" json:"metrics"`

	RegionID uuid.UUID `gorm:"type:uuid;index;not null" json:"region_id"`
	Region   Region    `gorm:"foreignKey:RegionID" json:"region,omitempty"`

	CheckID uuid.UUID `gorm:"type:uuid;index;not null" json:"check_id"`
	Check   Check     `gorm:"foreignKey:CheckID" json:"check,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// type Alert struct {
// 	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`

// 	ProjectID uuid.UUID `gorm:"type:uuid;index;not null" json:"project_id"`
// 	Project   Project   `gorm:"foreignKey:ProjectID" json:"project,omitempty"`

// 	CreatedAt time.Time      `json:"created_at"`
// 	UpdatedAt time.Time      `json:"updated_at"`
// 	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
// }

// type WebhookAttempt struct {
// 	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
// 	AlertID         *uuid.UUID     `gorm:"type:uuid;index" json:"alert_id,omitempty"`
// 	CheckID         uuid.UUID      `gorm:"type:uuid;index;not null" json:"check_id"`
// 	URL             string         `gorm:"type:text;not null" json:"url"`
// 	RequestBody     string         `gorm:"type:text" json:"request_body"`
// 	RequestHeaders  datatypes.JSON `gorm:"type:jsonb" json:"request_headers"`
// 	ResponseCode    *int           `gorm:"type:integer" json:"response_code,omitempty"`
// 	ResponseBody    *string        `gorm:"type:text" json:"response_body,omitempty"`
// 	ResponseHeaders datatypes.JSON `gorm:"type:jsonb" json:"response_headers"`
// 	Error           *string        `gorm:"type:text" json:"error,omitempty"`
// 	LatencyMs       *int64         `gorm:"type:bigint" json:"latency_ms,omitempty"`
// 	RetryNumber     int            `gorm:"default:0" json:"retry_number"`
// 	Timeout         bool           `gorm:"default:false" json:"timeout"`
// 	Status          string         `gorm:"type:varchar(20);not null" json:"status"` // success, failed
// 	CreatedAt       time.Time      `json:"created_at"`
// }
