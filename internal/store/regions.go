package store

import (
	"pulse/internal/models"
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
