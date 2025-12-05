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

func (s *Store) GetRegionByCode(code string) (*models.Region, error) {
	var region models.Region
	if err := s.db.Where("code = ?", code).First(&region).Error; err != nil {
		return nil, err
	}
	return &region, nil
}

func (s *Store) GetDueChecks(regionCode string) ([]models.Check, error) {
	var checks []models.Check
	now := time.Now()

	// Get region ID by code
	var region models.Region
	if err := s.db.Where("code = ?", regionCode).First(&region).Error; err != nil {
		return nil, err
	}

	// Get checks that are due and have this region enabled
	// Using a join to filter by region
	if err := s.db.
		Preload("Regions").
		Joins("JOIN check_regions ON checks.id = check_regions.check_id").
		Where("checks.is_enabled = ? AND check_regions.region_id = ? AND (checks.next_run_at IS NULL OR checks.next_run_at <= ?)", true, region.ID, now).
		Find(&checks).Error; err != nil {
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
