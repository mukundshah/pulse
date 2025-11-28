package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type Check struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name             string         `gorm:"not null"`
	URL              string         `gorm:"not null"`
	Method           string         `gorm:"not null;default:GET"`
	Headers          datatypes.JSON `gorm:"type:jsonb"`
	ExpectedStatus   int            `gorm:"default:200"`
	BodyContains     *string        `gorm:"type:text"`
	IntervalSeconds  int            `gorm:"not null"`
	AlertThreshold   int            `gorm:"not null;default:3"`
	ConsecutiveFails int            `gorm:"not null;default:0"`
	LastStatus       string         `gorm:"type:varchar(20);default:'unknown'"`
	LastRunAt        *time.Time
	NextRunAt        *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

type Alert struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CheckID   uuid.UUID `gorm:"type:uuid;index;not null"`
	Type      string    `gorm:"type:varchar(20);not null"`
	Payload   string    `gorm:"type:text;not null"`
	CreatedAt time.Time
	Check     Check     `gorm:"foreignKey:CheckID"`
	SentAt    time.Time `gorm:"not null"`
}
