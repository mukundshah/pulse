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

func (s *Store) GetCheckRunsByCheck(checkID uuid.UUID, limit int, after, before *uuid.UUID) ([]models.CheckRun, error) {
	var runs []models.CheckRun
	query := s.db.Preload("Region").Omit("Check").Where("check_id = ?", checkID)

	// Handle cursor-based pagination
	// We sort by created_at DESC, id DESC (latest to oldest)
	if after != nil {
		// Get the cursor item to find its created_at
		var cursorRun models.CheckRun
		if err := s.db.First(&cursorRun, "id = ?", *after).Error; err == nil {
			// Get items after the cursor (older items in DESC order)
			// (created_at < cursor.created_at) OR (created_at = cursor.created_at AND id < cursor.id)
			query = query.Where(
				"(created_at < ?) OR (created_at = ? AND id < ?)",
				cursorRun.CreatedAt, cursorRun.CreatedAt, *after,
			)
		}
	} else if before != nil {
		// Get the cursor item to find its created_at
		var cursorRun models.CheckRun
		if err := s.db.First(&cursorRun, "id = ?", *before).Error; err == nil {
			// Get items before the cursor (newer items in DESC order)
			// (created_at > cursor.created_at) OR (created_at = cursor.created_at AND id > cursor.id)
			query = query.Where(
				"(created_at > ?) OR (created_at = ? AND id > ?)",
				cursorRun.CreatedAt, cursorRun.CreatedAt, *before,
			)
		}
	}

	// Always sort from latest to oldest
	query = query.Order("created_at DESC, id DESC")

	if limit > 0 {
		// Fetch one extra to determine if there's a next page
		query = query.Limit(limit + 1)
	}

	if err := query.Find(&runs).Error; err != nil {
		return nil, err
	}
	return runs, nil
}
