package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Region struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Code      string         `gorm:"not null;unique" json:"code"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
