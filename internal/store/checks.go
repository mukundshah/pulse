package store

import (
	"time"

	"pulse/internal/models"

	"github.com/google/uuid"
)

func (s *Store) CreateCheck(check *models.Check) error {
	return s.db.Create(check).Error
}

func (s *Store) GetCheck(id uuid.UUID) (*models.Check, error) {
	var check models.Check
	if err := s.db.Preload("Project").Preload("Tags").Preload("Regions").First(&check, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &check, nil
}

func (s *Store) GetChecksByProject(projectID uuid.UUID) ([]models.Check, error) {
	var checks []models.Check
	if err := s.db.Preload("Tags").Preload("Regions").Where("project_id = ?", projectID).Find(&checks).Error; err != nil {
		return nil, err
	}
	return checks, nil
}

func (s *Store) UpdateCheck(check *models.Check) error {
	return s.db.Save(check).Error
}

func (s *Store) DeleteCheck(id uuid.UUID) error {
	return s.db.Delete(&models.Check{}, "id = ?", id).Error
}

func (s *Store) GetDueChecks() ([]models.Check, error) {
	var checks []models.Check
	now := time.Now()
	if err := s.db.Where("is_enabled = ? AND (next_run_at IS NULL OR next_run_at <= ?)", true, now).Find(&checks).Error; err != nil {
		return nil, err
	}
	return checks, nil
}

func (s *Store) UpdateCheckStatus(checkID uuid.UUID, nextRun time.Time, consecutiveFails int, lastStatus string) error {
	now := time.Now()
	return s.db.Model(&models.Check{}).Where("id = ?", checkID).Updates(map[string]interface{}{
		"last_run_at":       now,
		"next_run_at":       nextRun,
		"consecutive_fails": consecutiveFails,
		"last_status":       lastStatus,
	}).Error
}
