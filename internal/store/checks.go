package store

import (
	"pulse/internal/models"
	"pulse/internal/redis"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	db    *gorm.DB
	redis *redis.Client
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func NewWithCache(db *gorm.DB, r *redis.Client) *Store {
	return &Store{db: db, redis: r}
}

func (s *Store) GetDueChecks() ([]*models.Check, error) {
	var checks []*models.Check
	now := time.Now()
	err := s.db.Where("next_run_at <= ? OR next_run_at IS NULL", now).Find(&checks).Error
	return checks, err
}

func (s *Store) GetCheck(id uuid.UUID) (*models.Check, error) {
	if s.redis != nil {
		cached, err := s.redis.GetCheck(id)
		if err == nil && cached != nil {
			return cached, nil
		}
		// Cache miss or error, continue to DB
	}

	// Fallback to database
	var check models.Check
	err := s.db.First(&check, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	// Cache the result if Redis is available
	if s.redis != nil {
		// Cache for 5 minutes
		if err := s.redis.SetCheck(&check, 5*time.Minute); err != nil {
			// Log but don't fail on cache errors
			// log.Printf("Failed to cache check: %v", err)
		}
	}

	return &check, nil
}

func (s *Store) GetAllChecks() ([]*models.Check, error) {
	var checks []*models.Check
	err := s.db.Find(&checks).Error
	return checks, err
}

func (s *Store) CreateCheck(check *models.Check) error {
	if check.ID == uuid.Nil {
		check.ID = uuid.New()
	}

	err := s.db.Create(check).Error
	if err != nil {
		return err
	}

	// Invalidate cache
	if s.redis != nil {
		_ = s.redis.DeleteCheck(check.ID)
		// Cache the newly created check
		_ = s.redis.SetCheck(check, 5*time.Minute)
	}

	return nil
}

func (s *Store) UpdateCheckStatus(id uuid.UUID, nextRun time.Time, consecutiveFails int, lastStatus string) error {
	now := time.Now()
	err := s.db.Model(&models.Check{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"next_run_at":       nextRun,
			"last_run_at":       now,
			"consecutive_fails": consecutiveFails,
			"last_status":       lastStatus,
		}).Error
	if err != nil {
		return err
	}

	if s.redis != nil {
		_ = s.redis.DeleteCheck(id)
	}

	return nil
}
