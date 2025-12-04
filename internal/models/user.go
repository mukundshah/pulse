package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuidv4()" json:"id"`
	Name          string         `gorm:"not null" json:"name"`
	Email         string         `gorm:"not null;unique" json:"email"`
	PasswordHash  string         `gorm:"not null" json:"-"`
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	LastLogin     *time.Time     `json:"last_login,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
