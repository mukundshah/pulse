package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Alert struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`

	Status CheckRunStatus `gorm:"type:varchar(20);not null" json:"status"`

	RunID uuid.UUID `gorm:"type:uuid;index;not null" json:"run_id"`
	Run   CheckRun  `gorm:"foreignKey:RunID" json:"run,omitempty"`

	RegionID uuid.UUID `gorm:"type:uuid;index;not null" json:"region_id"`
	Region   Region    `gorm:"foreignKey:RegionID" json:"region,omitempty"`

	ProjectID uuid.UUID `gorm:"type:uuid;index;not null" json:"project_id"`
	Project   Project   `gorm:"foreignKey:ProjectID" json:"project,omitempty"`

	CheckID uuid.UUID `gorm:"type:uuid;index;not null" json:"check_id"`
	Check   Check     `gorm:"foreignKey:CheckID" json:"check,omitempty"`

	CreatedAt time.Time      `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`
}
