package store

import (
	"time"

	"pulse/internal/models"

	"github.com/google/uuid"
)

func (s *Store) CreateAlert(alert *models.Alert) error {
	return s.db.Create(alert).Error
}

func (s *Store) GetAlertsByCheck(checkID uuid.UUID, limit int, after, before *uuid.UUID, startTime, endTime *time.Time) ([]models.Alert, error) {
	var alerts []models.Alert
	query := s.db.Preload("Region").Preload("Run").Omit("Check").Where("check_id = ?", checkID)

	// Apply date range filter if provided
	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	// Handle cursor-based pagination
	// We sort by created_at DESC, id DESC (latest to oldest)
	if after != nil {
		// Get the cursor item to find its created_at
		var cursorAlert models.Alert
		if err := s.db.First(&cursorAlert, "id = ?", *after).Error; err == nil {
			// Get items after the cursor (older items in DESC order)
			// (created_at < cursor.created_at) OR (created_at = cursor.created_at AND id < cursor.id)
			query = query.Where(
				"(created_at < ?) OR (created_at = ? AND id < ?)",
				cursorAlert.CreatedAt, cursorAlert.CreatedAt, *after,
			)
		}
	} else if before != nil {
		// Get the cursor item to find its created_at
		var cursorAlert models.Alert
		if err := s.db.First(&cursorAlert, "id = ?", *before).Error; err == nil {
			// Get items before the cursor (newer items in DESC order)
			// (created_at > cursor.created_at) OR (created_at = cursor.created_at AND id > cursor.id)
			query = query.Where(
				"(created_at > ?) OR (created_at = ? AND id > ?)",
				cursorAlert.CreatedAt, cursorAlert.CreatedAt, *before,
			)
		}
	}

	// Always sort from latest to oldest
	query = query.Order("created_at DESC, id DESC")

	if limit > 0 {
		// Fetch one extra to determine if there's a next page
		query = query.Limit(limit + 1)
	}

	if err := query.Find(&alerts).Error; err != nil {
		return nil, err
	}

	return alerts, nil
}
