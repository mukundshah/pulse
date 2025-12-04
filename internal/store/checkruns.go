package store

import (
	"pulse/internal/models"

	"github.com/google/uuid"
)

func (s *Store) CreateCheckRun(run *models.CheckRun) error {
	return s.db.Create(run).Error
}

func (s *Store) GetCheckRun(id uuid.UUID) (*models.CheckRun, error) {
	var run models.CheckRun
	if err := s.db.Preload("Check").Preload("Region").First(&run, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &run, nil
}

func (s *Store) GetCheckRunsByCheck(checkID uuid.UUID, limit int) ([]models.CheckRun, error) {
	var runs []models.CheckRun
	query := s.db.Preload("Region").Where("check_id = ?", checkID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&runs).Error; err != nil {
		return nil, err
	}
	return runs, nil
}
