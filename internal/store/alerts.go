package store

import (
	"github.com/google/uuid"

	"pulse/internal/models"
)

func (s *Store) CreateAlert(alert *models.Alert) error {
	if alert.ID == uuid.Nil {
		alert.ID = uuid.New()
	}
	if alert.SentAt.IsZero() {
		alert.SentAt = alert.CreatedAt
	}
	return s.db.Create(alert).Error
}

func (s *Store) GetAlerts(checkID uuid.UUID, limit int) ([]*models.Alert, error) {
	var alerts []*models.Alert
	err := s.db.Where("check_id = ?", checkID).
		Omit("Check").
		Order("sent_at DESC").
		Limit(limit).
		Find(&alerts).Error
	return alerts, err
}
