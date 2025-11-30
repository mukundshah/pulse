package store

import (
	"pulse/internal/models"

	"github.com/google/uuid"
)

func (s *Store) CreateWebhookAttempt(attempt *models.WebhookAttempt) error {
	if attempt.ID == uuid.Nil {
		attempt.ID = uuid.New()
	}
	return s.db.Create(attempt).Error
}

func (s *Store) GetWebhookAttempts(checkID uuid.UUID, limit int) ([]*models.WebhookAttempt, error) {
	var attempts []*models.WebhookAttempt
	err := s.db.Where("check_id = ?", checkID).
		Order("created_at DESC").
		Limit(limit).
		Find(&attempts).Error
	return attempts, err
}

func (s *Store) GetWebhookAttemptsByAlert(alertID uuid.UUID) ([]*models.WebhookAttempt, error) {
	var attempts []*models.WebhookAttempt
	err := s.db.Where("alert_id = ?", alertID).
		Order("created_at ASC").
		Find(&attempts).Error
	return attempts, err
}
