package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	Name string    `gorm:"not null" json:"name"`

	CreatedAt time.Time      `gorm:"type:timestamptz" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`
}

type ProjectMember struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv4()" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid;index;not null" json:"project_id"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	Role      string    `gorm:"type:varchar(20);not null" json:"role"` // admin, member, viewer

	CreatedAt time.Time      `gorm:"type:timestamptz" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`

	Project Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type ProjectInvitation struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv4()" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid;index;not null" json:"project_id"`
	Email     string    `gorm:"not null" json:"email"`
	Token     string    `gorm:"not null" json:"token"`

	CreatedAt time.Time      `gorm:"type:timestamptz" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`
}
