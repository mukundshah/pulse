package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type FailureReason string

const (
	// Network
	FailureDNS                FailureReason = "dns_error"
	FailureTCP                FailureReason = "tcp_error"
	FailureTLS                FailureReason = "tls_error"
	FailureConnectionTimeout  FailureReason = "connection_timeout"
	FailureConnectionRefused  FailureReason = "connection_refused"
	FailureNetworkUnreachable FailureReason = "network_unreachable"

	// Timeouts
	FailureRequestTimeout  FailureReason = "request_timeout"
	FailureTTFBTimeout     FailureReason = "ttfb_timeout"
	FailureDownloadTimeout FailureReason = "download_timeout"
	FailureRunTimeout      FailureReason = "run_timeout"

	// HTTP
	FailureHTTP4xx          FailureReason = "http_4xx"
	FailureHTTP5xx          FailureReason = "http_5xx"
	FailureInvalidHTTP      FailureReason = "invalid_http_response"
	FailureUnexpectedStatus FailureReason = "unexpected_status_code"

	// Assertions
	FailureAssertionFailed  FailureReason = "assertion_failed"
	FailureContentMismatch  FailureReason = "content_mismatch"
	FailureHeaderMismatch   FailureReason = "header_mismatch"
	FailureSchemaValidation FailureReason = "schema_validation_failed"

	// Browser / Playwright
	FailureBrowserLaunch   FailureReason = "browser_launch_failed"
	FailureNavigation      FailureReason = "navigation_failed"
	FailureScript          FailureReason = "script_error"
	FailureElementNotFound FailureReason = "element_not_found"
	FailurePageCrash       FailureReason = "page_crash"

	// Orchestration
	FailureMaxRetries FailureReason = "max_retries_exceeded"
	FailureDependency FailureReason = "dependency_failed"

	// Internal
	FailureAgent         FailureReason = "agent_error"
	FailureSerialization FailureReason = "serialization_error"
	FailureUnknown       FailureReason = "unknown_error"
)

type CheckRun struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`

	Status             CheckRunStatus `gorm:"type:varchar(20);default:'unknown'" json:"status"`
	FailureReason      *FailureReason `gorm:"type:varchar(50)" json:"failure_reason,omitempty"`
	ResponseStatusCode *int32         `gorm:"type:integer" json:"response_status_code,omitempty"`

	RunStartedAt time.Time `gorm:"type:timestamptz;not null" json:"run_started_at"`
	RunEndedAt   time.Time `gorm:"type:timestamptz" json:"run_ended_at"`

	RequestStartedAt time.Time `gorm:"type:timestamptz;not null" json:"request_started_at"`
	FirstByteAt      time.Time `gorm:"type:timestamptz" json:"first_byte_at"`
	ResponseEndedAt  time.Time `gorm:"type:timestamptz" json:"response_ended_at"`

	ConnectionReused  bool  `gorm:"type:boolean;default:false" json:"connection_reused"`
	ResponseSizeBytes int64 `gorm:"type:bigint" json:"response_size_bytes"`

	AssertionResults datatypes.JSON `gorm:"type:jsonb" json:"assertion_results"`
	PlaywrightReport datatypes.JSON `gorm:"type:jsonb" json:"playwright_report,omitempty"`
	NetworkTimings   datatypes.JSON `gorm:"type:jsonb" json:"network_timings"`

	RegionID uuid.UUID `gorm:"type:uuid;index;not null" json:"region_id"`
	Region   Region    `gorm:"foreignKey:RegionID" json:"region,omitempty"`

	CheckID uuid.UUID `gorm:"type:uuid;index;not null" json:"check_id"`
	Check   *Check    `gorm:"foreignKey:CheckID" json:"check,omitempty"`

	CreatedAt time.Time      `gorm:"type:timestamptz" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`
}
