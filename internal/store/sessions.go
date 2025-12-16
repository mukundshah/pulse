package store

import (
	"time"

	"pulse/internal/models"

	"github.com/google/uuid"
)

// CreateSession creates a new session for a user
func (s *Store) CreateSession(session *models.Session) error {
	if session.ID == uuid.Nil {
		session.ID = uuid.New()
	}
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now().UTC()
	}
	if session.UpdatedAt.IsZero() {
		session.UpdatedAt = time.Now().UTC()
	}
	if session.LastActivity.IsZero() {
		session.LastActivity = time.Now().UTC()
	}

	err := s.db.Create(session).Error
	if err != nil {
		return err
	}

	// Optionally cache in Redis if available
	if s.redis != nil {
		ttl := time.Until(session.ExpiresAt)
		if ttl > 0 {
			// Store session metadata in Redis for fast lookups
			if err := s.redis.SetSession(session.JTI, ttl); err != nil {
				// Log error but don't fail session creation
				// This is a cache, not critical
			}
		}
	}

	return nil
}

// GetSessionByJTI retrieves a session by its JWT ID (JTI)
func (s *Store) GetSessionByJTI(jti string) (*models.Session, error) {
	// Try Redis first if available
	if s.redis != nil {
		val, err := s.redis.GetSession(jti)
		if err != nil {
			// Not in cache, might be invalidated, but check DB anyway
		} else if val == "" {
			// Explicitly invalidated in cache
			return nil, ErrNotFound
		}
	}

	var session models.Session
	err := s.db.Where("jti = ? AND is_active = ? AND expires_at > ?", jti, true, time.Now().UTC()).First(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// GetSessionByID retrieves a session by its ID
func (s *Store) GetSessionByID(sessionID uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := s.db.Where("id = ?", sessionID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetSessionsByUserID retrieves all active sessions for a user
func (s *Store) GetSessionsByUserID(userID uuid.UUID) ([]*models.Session, error) {
	var sessions []*models.Session
	err := s.db.Where("user_id = ? AND is_active = ? AND expires_at > ?", userID, true, time.Now().UTC()).
		Order("last_activity DESC").
		Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// InvalidateSession marks a session as inactive by JTI
func (s *Store) InvalidateSession(jti string) error {
	// Update database
	err := s.db.Model(&models.Session{}).
		Where("jti = ?", jti).
		Updates(map[string]interface{}{
			"is_active":  false,
			"updated_at": time.Now().UTC(),
		}).Error
	if err != nil {
		return err
	}

	// Remove from Redis cache if available
	if s.redis != nil {
		s.redis.DeleteSession(jti)
	}

	return nil
}

// InvalidateAllUserSessions invalidates all sessions for a user (useful for password reset, security breach, etc.)
func (s *Store) InvalidateAllUserSessions(userID uuid.UUID) error {
	// Get all active sessions first to remove from Redis
	if s.redis != nil {
		sessions, err := s.GetSessionsByUserID(userID)
		if err == nil {
			for _, session := range sessions {
				s.redis.DeleteSession(session.JTI)
			}
		}
	}

	// Update database
	err := s.db.Model(&models.Session{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Updates(map[string]interface{}{
			"is_active":  false,
			"updated_at": time.Now().UTC(),
		}).Error

	return err
}

// UpdateSessionActivity updates the last activity timestamp for a session
func (s *Store) UpdateSessionActivity(jti string) error {
	return s.db.Model(&models.Session{}).
		Where("jti = ?", jti).
		Update("last_activity", time.Now().UTC()).Error
}

// CleanupExpiredSessions removes expired sessions from the database
// This should be run periodically (e.g., via a cron job)
func (s *Store) CleanupExpiredSessions() error {
	// Delete expired sessions
	result := s.db.Where("expires_at < ?", time.Now().UTC()).Delete(&models.Session{})
	return result.Error
}

// ErrNotFound is returned when a session is not found
var ErrNotFound = &NotFoundError{}

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "session not found"
}
