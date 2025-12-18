package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pulse/internal/middleware"
	"pulse/internal/store"
)

type AlertHandler struct {
	store *store.Store
}

func NewAlertHandler(s *store.Store) *AlertHandler {
	return &AlertHandler{store: s}
}

// ListAlerts handles GET /projects/:projectId/checks/:checkId/alerts
// Returns the last 10 alerts for a check
func (h *AlertHandler) ListAlerts(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	isMember, err := h.store.IsProjectMember(projectID, userID)
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	checkID, err := uuid.Parse(c.Param("checkId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check ID"})
		return
	}

	// Fetch the last 10 alerts (limit = 10, no pagination for now)
	alerts, err := h.store.GetAlertsByCheck(checkID, 10, nil, nil, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list alerts"})
		return
	}

	response := gin.H{
		"data": alerts,
	}

	c.JSON(http.StatusOK, response)
}
