package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Session represents an active user session
type Session struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuidv4()" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	JTI          string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"jti"` // JWT ID
	UserAgent    string         `gorm:"type:text" json:"user_agent,omitempty"`
	IPAddress    string         `gorm:"type:varchar(45)" json:"ip_address,omitempty"` // IPv6 max length
	IsActive     bool           `gorm:"default:true;index" json:"is_active"`
	ExpiresAt    time.Time      `gorm:"index;not null" json:"expires_at"`
	LastActivity time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"last_activity"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
