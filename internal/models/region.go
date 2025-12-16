package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Region struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	Name string    `gorm:"not null" json:"name"`
	Code string    `gorm:"not null;unique" json:"code"`
	Flag string    `gorm:"type:varchar(10)" json:"flag"`

	CreatedAt time.Time      `gorm:"type:timestamptz" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`
}
