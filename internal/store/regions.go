package store

import (
	"pulse/internal/models"

	"github.com/google/uuid"
)

// ListRegions returns all system regions (read-only)
func (s *Store) ListRegions() ([]models.Region, error) {
	var regions []models.Region
	if err := s.db.Find(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}

func (s *Store) GetRegionByCode(code string) (*models.Region, error) {
	var region models.Region
	if err := s.db.Where("code = ?", code).First(&region).Error; err != nil {
		return nil, err
	}
	return &region, nil
}

func (s *Store) AddRegionToCheck(checkID uuid.UUID, regionID uuid.UUID) error {
	var check models.Check
	if err := s.db.First(&check, "id = ?", checkID).Error; err != nil {
		return err
	}
	var region models.Region
	if err := s.db.First(&region, "id = ?", regionID).Error; err != nil {
		return err
	}
	return s.db.Model(&check).Association("Regions").Append(&region)
}

func (s *Store) RemoveRegionFromCheck(checkID uuid.UUID, regionID uuid.UUID) error {
	var check models.Check
	if err := s.db.First(&check, "id = ?", checkID).Error; err != nil {
		return err
	}
	var region models.Region
	if err := s.db.First(&region, "id = ?", regionID).Error; err != nil {
		return err
	}
	return s.db.Model(&check).Association("Regions").Delete(&region)
}
