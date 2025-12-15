package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Check struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	Name string    `gorm:"not null" json:"name"`

	IsEnabled  bool `gorm:"default:true" json:"is_enabled"`
	IsMuted    bool `gorm:"default:false" json:"is_muted"`
	ShouldFail bool `gorm:"default:false" json:"should_fail"`

	Type CheckType `gorm:"not null;default:http" json:"type"`

	Host        string         `gorm:"not null" json:"host"` // could be a domain or an IP address
	Port        int            `gorm:"default:80" json:"port"`
	Secure      bool           `gorm:"default:false" json:"secure"`
	Method      string         `gorm:"not null;default:GET" json:"method"`
	Path        string         `gorm:"not null;default:/" json:"path"`
	QueryParams datatypes.JSON `gorm:"type:jsonb" json:"query_params"`
	Headers     datatypes.JSON `gorm:"type:jsonb" json:"headers"`
	Body        datatypes.JSON `gorm:"type:jsonb" json:"body"`
	IPVersion   IPVersionType  `gorm:"not null;default:ipv4" json:"ip_version"`

	SkipSSLVerification bool `gorm:"default:false" json:"skip_ssl_verification"`
	FollowRedirects     bool `gorm:"default:true" json:"follow_redirects"`

	PlaywrightScript *string        `gorm:"type:text" json:"playwright_script,omitempty"`
	Assertions       datatypes.JSON `gorm:"type:jsonb" json:"assertions"`

	PreScript  *string `gorm:"type:text" json:"pre_script,omitempty"`
	PostScript *string `gorm:"type:text" json:"post_script,omitempty"`

	Interval string `gorm:"not null;default:'10m'" json:"interval"`

	DegradedThreshold     int      `gorm:"not null" json:"degraded_threshold"`
	DegradedThresholdUnit UnitType `gorm:"type:varchar(2);default:'ms'" json:"degraded_threshold_unit"`
	FailedThreshold       int      `gorm:"not null" json:"failed_threshold"`
	FailedThresholdUnit   UnitType `gorm:"type:varchar(2);default:'ms'" json:"failed_threshold_unit"`

	Retries             RetryType        `gorm:"type:varchar(20);default:'none'" json:"retries"`
	RetriesCount        *int             `json:"retries_count,omitempty"`
	RetriesDelay        *int             `json:"retries_delay,omitempty"`
	RetriesDelayUnit    *UnitType        `json:"retries_delay_unit,omitempty"`
	RetriesFactor       *float64         `json:"retries_factor,omitempty"`
	RetriesJitter       *RetryJitterType `json:"retries_jitter,omitempty"`
	RetriesJitterFactor *float64         `json:"retries_jitter_factor,omitempty"`
	RetriesMaxDelay     *int             `json:"retries_max_delay,omitempty"`
	RetriesMaxDelayUnit *UnitType        `json:"retries_max_delay_unit,omitempty"`
	RetriesTimeout      *int             `json:"retries_timeout,omitempty"`
	RetriesTimeoutUnit  *UnitType        `json:"retries_timeout_unit,omitempty"`

	DNSRecordType       *DNSRecordType           `json:"dns_record_type,omitempty"`
	DNSResolver         *string                  `json:"dns_resolver,omitempty"`
	DNSResolverPort     *int                     `json:"dns_resolver_port,omitempty"`
	DNSResolverProtocol *DNSResolverProtocolType `json:"dns_resolver_protocol,omitempty"`

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

func (c *Check) FailedThresholdDuration() time.Duration {
	duration, err := time.ParseDuration(fmt.Sprintf("%d%s", c.FailedThreshold, c.FailedThresholdUnit))
	if err != nil {
		return 0
	}
	return duration
}

func (c *Check) DegradedThresholdDuration() time.Duration {
	duration, err := time.ParseDuration(fmt.Sprintf("%d%s", c.DegradedThreshold, c.DegradedThresholdUnit))
	if err != nil {
		return 0
	}
	return duration
}

func (c *Check) RetriesDelayDuration() time.Duration {
	if c.RetriesDelay == nil || c.RetriesDelayUnit == nil {
		return 0
	}
	duration, err := time.ParseDuration(fmt.Sprintf("%d%s", *c.RetriesDelay, *c.RetriesDelayUnit))
	if err != nil {
		return 0
	}
	return duration
}

func (c *Check) RetriesMaxDelayDuration() time.Duration {
	if c.RetriesMaxDelay == nil || c.RetriesMaxDelayUnit == nil {
		return 0
	}
	duration, err := time.ParseDuration(fmt.Sprintf("%d%s", *c.RetriesMaxDelay, *c.RetriesMaxDelayUnit))
	if err != nil {
		return 0
	}
	return duration
}

func (c *Check) RetriesTimeoutDuration() time.Duration {
	if c.RetriesTimeout == nil || c.RetriesTimeoutUnit == nil {
		return 0
	}
	duration, err := time.ParseDuration(fmt.Sprintf("%d%s", *c.RetriesTimeout, *c.RetriesTimeoutUnit))
	if err != nil {
		return 0
	}
	return duration
}

func (c *Check) IntervalDuration() time.Duration {
	duration, err := time.ParseDuration(c.Interval)
	if err != nil {
		return 0
	}
	return duration
}
