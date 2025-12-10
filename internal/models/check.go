package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type IPVersionType string

const (
	IPVersionTypeIPv4 IPVersionType = "ipv4"
	IPVersionTypeIPv6 IPVersionType = "ipv6"
)

type Check struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	Name string    `gorm:"not null" json:"name"`

	IsEnabled bool `gorm:"default:true" json:"is_enabled"`
	IsMuted   bool `gorm:"default:false" json:"is_muted"`

	Type CheckType `gorm:"not null;default:http" json:"type"`

	Host        string         `gorm:"not null" json:"host"` // could be a domain or an IP address
	Port        int            `gorm:"default:80" json:"port"`
	Secure      bool           `gorm:"default:false" json:"secure"`
	Method      string         `gorm:"not null;default:GET" json:"method"`
	Path        string         `gorm:"not null;default:/" json:"path"`
	Headers     datatypes.JSON `gorm:"type:jsonb" json:"headers"`
	QueryParams datatypes.JSON `gorm:"type:jsonb" json:"query_params"`
	Body        datatypes.JSON `gorm:"type:jsonb" json:"body"`
	IPVersion   IPVersionType  `gorm:"not null;default:ipv4" json:"ip_version"`

	SSLVerification bool `gorm:"default:true" json:"ssl_verification"`
	FollowRedirects bool `gorm:"default:true" json:"follow_redirects"`

	PlaywrightScript *string        `gorm:"type:text" json:"playwright_script,omitempty"`
	Assertions       datatypes.JSON `gorm:"type:jsonb" json:"assertions"`
	ExpectedStatus   int            `gorm:"default:200" json:"expected_status"`
	ShouldFail       bool           `gorm:"default:false" json:"should_fail"`

	PreScript  *string `gorm:"type:text" json:"pre_script,omitempty"`
	PostScript *string `gorm:"type:text" json:"post_script,omitempty"`

	TimeoutMs       int `gorm:"default:10000" json:"timeout_ms"`
	IntervalSeconds int `gorm:"not null" json:"interval_seconds"`
	AlertThreshold  int `gorm:"not null;default:3" json:"alert_threshold"`

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
