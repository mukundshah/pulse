package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	Name string    `gorm:"not null" json:"name"`

	CreatedAt time.Time      `gorm:"type:timestamptz" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`

	ProjectID uuid.UUID `gorm:"type:uuid;index;not null" json:"project_id"`
	Project   Project   `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}
