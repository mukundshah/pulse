package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pulse/internal/middleware"
	"pulse/internal/store"
)

type SessionHandler struct {
	store *store.Store
}

func NewSessionHandler(s *store.Store) *SessionHandler {
	return &SessionHandler{store: s}
}

// SessionResponse represents a session in API responses
type SessionResponse struct {
	ID           uuid.UUID `json:"id"`
	UserAgent    string    `json:"user_agent,omitempty"`
	IPAddress    string    `json:"ip_address,omitempty"`
	IsActive     bool      `json:"is_active"`
	IsCurrent    bool      `json:"is_current"` // True if this is the current session
	ExpiresAt    string    `json:"expires_at"`
	LastActivity string    `json:"last_activity"`
	CreatedAt    string    `json:"created_at"`
}

// ListSessions returns all active sessions for the current user
// GET /sessions
func (h *SessionHandler) ListSessions(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Get current session ID from context
	currentSessionID, _ := c.Get("session_id")
	currentSessionUUID, _ := currentSessionID.(uuid.UUID)

	// Get all sessions for user
	sessions, err := h.store.GetSessionsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve sessions"})
		return
	}

	// Convert to response format
	response := make([]SessionResponse, len(sessions))
	for i, session := range sessions {
		response[i] = SessionResponse{
			ID:           session.ID,
			UserAgent:    session.UserAgent,
			IPAddress:    session.IPAddress,
			IsActive:     session.IsActive,
			IsCurrent:    session.ID == currentSessionUUID,
			ExpiresAt:    session.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
			LastActivity: session.LastActivity.Format("2006-01-02T15:04:05Z07:00"),
			CreatedAt:    session.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"sessions": response})
}

// RevokeCurrentSession revokes the current session
// DELETE /sessions/current
func (h *SessionHandler) RevokeCurrentSession(c *gin.Context) {
	jti, exists := c.Get("jti")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	jtiStr, ok := jti.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	if err := h.store.InvalidateSession(jtiStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session revoked successfully"})
}

// RevokeSession revokes a specific session by ID
// DELETE /sessions/:sessionId
func (h *SessionHandler) RevokeSession(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	sessionID, err := uuid.Parse(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	// Get session to verify it belongs to the user
	session, err := h.store.GetSessionByID(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	// Verify session belongs to current user
	if session.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "session does not belong to current user"})
		return
	}

	// Revoke session by JTI
	if err := h.store.InvalidateSession(session.JTI); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session revoked successfully"})
}

// RevokeSessionsRequest represents a request to revoke multiple sessions
type RevokeSessionsRequest struct {
	SessionIDs []uuid.UUID `json:"session_ids" binding:"required,min=1"`
}

// RevokeSessions revokes multiple sessions by their IDs
// DELETE /sessions/batch
func (h *SessionHandler) RevokeSessions(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req RevokeSessionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current session JTI to prevent revoking current session
	currentJTI, _ := c.Get("jti")
	currentJTIStr, _ := currentJTI.(string)

	revokedCount := 0
	errors := make([]string, 0)

	for _, sessionID := range req.SessionIDs {
		// Get session to verify it belongs to the user
		session, err := h.store.GetSessionByID(sessionID)
		if err != nil {
			errors = append(errors, "session not found: "+sessionID.String())
			continue
		}

		// Verify session belongs to current user
		if session.UserID != userID {
			errors = append(errors, "session does not belong to current user: "+sessionID.String())
			continue
		}

		// Prevent revoking current session (user should use /sessions/current)
		if session.JTI == currentJTIStr {
			errors = append(errors, "cannot revoke current session: "+sessionID.String())
			continue
		}

		// Revoke session
		if err := h.store.InvalidateSession(session.JTI); err != nil {
			errors = append(errors, "failed to revoke session: "+sessionID.String())
			continue
		}

		revokedCount++
	}

	response := gin.H{
		"message":      "batch revocation completed",
		"revoked_count": revokedCount,
		"total_count":   len(req.SessionIDs),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	c.JSON(http.StatusOK, response)
}

// RevokeAllSessions revokes all sessions for the current user except the current one
// DELETE /sessions/all
func (h *SessionHandler) RevokeAllSessions(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Get current session JTI to exclude it from revocation
	currentJTI, _ := c.Get("jti")
	currentJTIStr, _ := currentJTI.(string)

	// Get all sessions for user
	sessions, err := h.store.GetSessionsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve sessions"})
		return
	}

	revokedCount := 0
	for _, session := range sessions {
		// Skip current session
		if session.JTI == currentJTIStr {
			continue
		}

		if err := h.store.InvalidateSession(session.JTI); err != nil {
			// Log error but continue with other sessions
			continue
		}
		revokedCount++
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "all other sessions revoked successfully",
		"revoked_count": revokedCount,
	})
}
