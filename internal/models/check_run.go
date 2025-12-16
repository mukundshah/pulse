package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CheckRun struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`

	Status         CheckRunStatus `gorm:"type:varchar(20);default:'unknown'" json:"status"`
	ResponseStatus int32          `gorm:"type:integer" json:"response_status"`
	TotalTimeMs    int            `gorm:"type:integer" json:"total_time_ms"`

	AssertionResults datatypes.JSON `gorm:"type:jsonb" json:"assertion_results"`
	PlaywrightReport datatypes.JSON `gorm:"type:jsonb" json:"playwright_report,omitempty"`
	NetworkTimings   datatypes.JSON `gorm:"type:jsonb" json:"network_timings"`

	RegionID uuid.UUID `gorm:"type:uuid;index;not null" json:"region_id"`
	Region   Region    `gorm:"foreignKey:RegionID" json:"region,omitempty"`

	CheckID uuid.UUID `gorm:"type:uuid;index;not null" json:"check_id"`
	Check   *Check    `gorm:"foreignKey:CheckID" json:"check,omitempty"`

	Remarks string `gorm:"type:text" json:"remarks"`

	CreatedAt time.Time      `gorm:"type:timestamptz" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`
}
