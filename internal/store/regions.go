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
